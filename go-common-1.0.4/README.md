# aoplatform common
# 控制台通用工具包
package
```shell
"github.com/eolinker/go-common"
```

## aolabel
定义
```golang

type Struct struct{

	Creator auto.Label `json:"creator" aolabel:"user"`

}
```
auto 赋值
```golang

v:= &Struct{
Creator: auto.UUID("uuid")
}

list:= make([]*Struct{},0)
...

auto.CompleteLabels(cxt,v)
auto.CompleteLabels(ctx,list)
auto.CompleteLabels(ctx,list...)
```


注册完成器

```golang
import (
    "github/eolinker/go-common/auto"
)

auto.RegisterService(name, handler)
```

