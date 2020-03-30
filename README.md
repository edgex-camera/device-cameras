# device-cameras

统一摄像头服务，整合不同类型摄像头

## 本每个edgex设备包含2个部分
1. 摄像头，必须包含
2. 控制协议，如onvif、sip等

## 已知设备类型
| 编号  | 设备名称 | 摄像头类型 | 控制协议 | 备注 |  已支持
| :--- | :-----: | :-------: | :-----: | :--: | :--: |
| 1 | normal-camera | usb/rtsp | - | 普通usb/ip摄像头　| √ |
| 2 | onvif-camera | rtsp | onvif | 球机onvif摄像头 | √ |
| 3 | dual-usb-camera | usb | - | 双usb摄像头 | ×　|
| 4 | simple-camera | usb/rtsp | - | 简化版摄像头，仅支持拍照 |　×　|

## 配置存储方式
```
|_ Driver
    |_ all_devices(json，key为device name, value为true)
    |_ "DeviceName1" (json, 该device专属配置)
        |_ DeviceName1.camera (json, 该device camera配置)
            |_ DeviceName1.camera."channelId1"
            |_ DeviceName1.camera."channelId2"
        |_ DeviceName1.control (json, 该device控制协议配置)
            |_ DeviceName1.onvif.config (json, onvif基础配置)
            |_ DeviceName1.onvif.presets (json, onvif预置位信息)
    |_ "DeviceName2"
        ...
```
