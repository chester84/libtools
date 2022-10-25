# libtools
```
common golang lib tools

1. try v1.0.0
2. try v1.0.1
3. try v1.0.2
4. try v1.0.3
```

## some functions specifications

### date tools
```
//golang origin date functions drive me crazy
//UnixMsec2Date and Date2UnixMsec are better ones O(∩_∩)O

showTime := libtools.UnixMsec2Date(1664182378999, "Y-m-d H:i:s")
fmt.Println("showTime:", showTime)
//showTime: 2022-09-26 16:52:58

ut := libtools.Date2UnixMsec("2022-09-26 16:52:58", "Y-m-d H:i:s")
fmt.Println("ut:", ut)
//ut: 1664182378000
```

### str,int64,int conversion
```
libtools.Int642Str(1664182378000)
libtools.Str2Int64("1664182378000")
libtools.AbsInt64(1)
libtools.AbsInt64(-1)
```

### Struct2Map
```
libtools.Struct2Map
libtools.Map2struct
```

### MoneyDisplay
```
libtools.MoneyDisplay
libtools.HumanMoney
```

### http request
```
libtools.SimpleHttpClient
```

### md5
```
libtools.Sha256()
libtools.Sha1()
libtools.Md5()
```
