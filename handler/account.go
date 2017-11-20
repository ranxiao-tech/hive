package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/encrypt"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
)

func UpdatePassword(c echo.Context) error {
	pwd := c.FormValue("pwd")
	pwdNew := c.FormValue("pwdNew")
	itf := c.Get("account")
	if itf == nil {
		return utils.AuthFailNull(c)
	}
	acc := global.ToInterfaceAccount(itf)
	if acc == nil {
		return utils.AuthFailNull(c)
	}
	sql := `SELECT id, password,salt FROM account WHERE id = ? AND status = ?`
	rows, _ := global.DB.Query(sql, acc.ID, enum.NORMAL)
	if len(rows) != 1 {
		return utils.Error(c, "帐号不存在或已禁用", nil)
	}
	userInfo := rows[0]
	pwdOld := convert.ToString(userInfo["password"])
	salt := convert.ToString(userInfo["salt"])
	if encrypt.Sha1Encode(pwd+salt) != pwdOld {
		return utils.Error(c, "原密码错误", nil)
	}

	sql = `UPDATE account SET password=? WHERE id = ? AND status = ?`
	x, err := global.DB.Update(sql, encrypt.Sha1Encode(pwdNew+salt), acc.ID, enum.NORMAL)
	if err != nil {
		return utils.ErrorNull(c, err.Error())
	}
	if x > 0 {
		return utils.Success(c, "修改密码成功", nil)
	} else {
		return utils.Error(c, "修改密码失败", nil)
	}
}