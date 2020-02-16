# device-cameras

统一摄像头服务，整合不同类型摄像头
[相关文档](http://jxdata.jiangxingai.com:10002/oo/r/535650527069912150)

## 配置存储方式
| Driver
    | all_devices(json，key为device name, value为true)
    | "DeviceName1" (各个device专属配置)
        | basic (json, 该device基础信息)
        | camera (json, 该device camera配置)
            | "channelId1"
            | "channelId2"
        | onvif (json, 该device onvif配置)
            | config (json, onvif基础配置)
            | presets (json, 预置位信息)
    | "DeviceName2"
        ...
