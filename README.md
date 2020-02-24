# device-cameras

统一摄像头服务，整合不同类型摄像头
[相关文档](http://jxdata.jiangxingai.com:10002/oo/r/535650527069912150)

## 配置存储方式
```
|_ Driver
    |_ all_devices(json，key为device name, value为true)
    |_ "DeviceName1" (json, 该device专属配置)
        |_ DeviceName1.camera (json, 该device camera配置)
            |_ DeviceName1.camera."channelId1"
            |_ DeviceName1.camera"channelId2"
        |_ DeviceName1.onvif (json, 该device onvif配置)
            |_ DeviceName1.onvif.config (json, onvif基础配置)
            |_ DeviceName1.onvvif.presets (json, 预置位信息)
    |_ "DeviceName2"
        ...
```
