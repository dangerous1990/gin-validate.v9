# gin-validate.v9
使用validate.v9 替换validate.v8
# 使用方法

## (一) 如何使用validate.v9标签进行参数校验

1. 参考官方文档 https://godoc.org/gopkg.in/go-playground/validator.v9 tag标签如何使用

2. validate.v9 默认使用validate，validate可以自定义tag, gin集成validate.v8的时候，tag使用了binding标签，为了兼容以前的代码，继续使用binding标签

## (二)  自定义验证器
1.定义一个方法，返回值是bool类型，参数是validator.FieldLevel
```
// 校验以“,”分隔的多个邮箱 eg. test1@163.com,test2@163.com
func ValidateMultiEmails(fl validator.FieldLevel) bool {
	return helper.IsValidMultiEmails(fl.Field().String())
}

```
2. 可以通过fl.Filed()获取字段的类型、值等信息，fl.Filed()是这个属性的反射值reflect.Value
3. 注册到***v8_to_v9.go*** lazyinit 方法中

v.validate.RegisterValidation("isValidMultiEmails", ValidateMultiEmails)

```
func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		glog.Infof("init gin validate")
		// 国际化
		localZH := local_zh.New()
		uni = ut.New(localZH, localZH)
		trans, _ := uni.GetTranslator("zh")

		v.validate = validator.New()
		v.trans = trans
		// validate.v9 tag 默认validate，兼容老代码
		v.validate.SetTagName("binding")
		// 汉化验证提示
		translations_zh.RegisterDefaultTranslations(v.validate, trans)
		// 自定义验证器 https://godoc.org/gopkg.in/go-playground/validator.v9
		v.validate.RegisterValidation("isValidMultiEmails", ValidateMultiEmails)
	})
}
```
4. 在struct中binding使用
```
		ContactsReceiveTaxfileMails string `form:"contacts_receive_taxfile_mails" json:"contacts_receive_taxfile_mails" binding:"required,isValidMultiEmails"`

```

## 汉化自定义验证器返回详细内容
1. 自定义的验证器默认会返回的英文的验证结果，需要对验证结果进行汉化
```
func ValidateMultiEmailsRegisterTranslationsFunc(ut ut.Translator) (err error) {
	if err = ut.Add("isValidMultiEmails", "{0}邮箱不合法，多个邮箱请以逗号(半角)分隔", false); err != nil {
		return
	}
	return

}
func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		glog.Warningf("警告: 翻译字段错误: %#v", fe)
		return fe.(error).Error()
	}
	return t
}

```
2. 注册到validate中

```
v.validate.RegisterTranslation("isValidMultiEmails", trans, ValidateMultiEmailsRegisterTranslationsFunc, translateFunc)
```
结果展示
```
{
    "error": "ContactsReceiveTaxfileMails邮箱不合法，多个邮箱请以逗号(半角)分隔",
}
```
3. 更多汉化请参考 gopkg.in/go-playground/validator.v9/translations/zh


# 常见的问题
1. 对于数字类型使用required，在传递0值的时候会报错

解决办法： 可以加上isdefault,这样就可以传递0值了
```
type struct form {
    id int 'json:"id" binding:"required,isdefault"'
}
```
2. 对于字符串类型，慎用len标签

解决办法：使用自定义验证器替代



