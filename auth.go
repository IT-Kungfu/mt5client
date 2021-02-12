package mt5client

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/text/encoding/unicode"
)

func (c *MT5Client) auth() error {
	authStart, body, err := c.authStartRequest()
	if err != nil {
		return err
	}

	c.log.Debugf("Auth start header: %s body: %s", authStart[:MT5HeaderLength], body)

	cmd, err := c.sendRequest(authStart)
	if err != nil {
		return err
	}

	if cmd.Params[MT5RetCode] != MT5RetCodeSuccess {
		return fmt.Errorf("auth retcode error: %s", cmd.Params[MT5RetCode])
	}

	c.log.Debugf("Auth start response: %+v", cmd)

	authAnswer, body, err := c.authAnswerRequest(cmd)
	if err != nil {
		return err
	}

	c.log.Debugf("Auth answer header: %s body: %s", authAnswer[:MT5HeaderLength], body)

	cmd, err = c.sendRequest(authAnswer)
	if err != nil {
		return err
	}

	if cmd.Params[MT5RetCode] != MT5RetCodeSuccess {
		return fmt.Errorf("auth retcode error: %s", cmd.Params[MT5RetCode])
	}

	c.log.Debugf("Auth answer response: %+v", cmd)

	return nil
}

func (c *MT5Client) authStartRequest() ([]byte, string, error) {
	body := fmt.Sprintf("%s|VERSION=%s|AGENT=%s|LOGIN=%s|TYPE=MANAGER|CRYPT_METHOD=NONE|%s",
		MT5CommandAuthStart, c.cfg.MT5APIVersion, c.cfg.MT5APIAgent, c.cfg.MT5Login, MT5PacketSeparator)
	req, err := makePacket(body, 0, 0)
	return req, body, err
}

func (c *MT5Client) authAnswerRequest(cmd *MT5Command) ([]byte, string, error) {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	encPassword, err := utf16.NewEncoder().String(c.cfg.MT5Password)
	if err != nil {
		return nil, "", err
	}

	tmpPasswordHash := md5.Sum([]byte(encPassword))
	tmpPasswordHash2 := md5.Sum(append(tmpPasswordHash[:], []byte("WebAPI")...))
	arrSrvRand, err := hex.DecodeString(cmd.Params["SRV_RAND"])
	if err != nil {
		return nil, "", err
	}
	randAnswer := fmt.Sprintf("%x", md5.Sum(append(tmpPasswordHash2[:], arrSrvRand...)))

	cliRand := fmt.Sprintf("%x", md5.Sum(makeRandomString()))
	body := fmt.Sprintf("%s|SRV_RAND_ANSWER=%s|CLI_RAND=%s|%s",
		MT5CommandAuthAnswer, randAnswer, cliRand, MT5PacketSeparator)

	req, err := makePacket(body, 0, 0)
	return req, body, err
}
