# device-cameras

统一摄像头服务，整合不同类型摄像头
[相关文档](http://jxdata.jiangxingai.com:10002/oo/r/535650527069912150)

## 配置存储方式
|_ Driver
    |_ all_devices(json，key为device name, value为true)
    |_ "DeviceName1" (各个device专属配置)
        |_ basic (json, 该device基础信息)
        |_ camera (json, 该device camera配置)
            |_ "channelId1"
            |_ "channelId2"
        |_ onvif (json, 该device onvif配置)
            |_ config (json, onvif基础配置)
            |_ presets (json, 预置位信息)
    |_ "DeviceName2"
        ...
