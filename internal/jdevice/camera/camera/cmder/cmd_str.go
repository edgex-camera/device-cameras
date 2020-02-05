package cmder

type cmdTemplate struct {
	base    string
	quality string
	capture string
	h264    string
	stream  string
	video   string
}

type processCmdTemplate struct {
	processor string
	rtsp      cmdTemplate
	webcam    cmdTemplate
}

var (
	ffmpeg = processCmdTemplate{
		processor: "ffmpeg",
		rtsp: cmdTemplate{
			base:    "-rtsp_transport tcp -i %s ",
			capture: "-r %d/1 -strftime 1 -y %s ",
			stream:  "-vcodec copy -an -f flv %s ",
			video: `-flags +global_header 
		-f stream_segment -segment_time %d -segment_format_options 
		movflags=+faststart -reset_timestamps 1 
		-vcodec copy -q:v 4 -an -r 24 -strftime 1 %s `,
		},
		webcam: cmdTemplate{
			base:    "-use_wallclock_as_timestamps 1 -f v4l2 -vcodec mjpeg -s %[1]dx%[2]d -i %[3]s ",
			capture: "-r %d/1 -strftime 1 -y %s ",
			stream:  "-vcodec h264 -an -f flv %s ",
			video: `-flags +global_header 
			-f stream_segment -segment_time %d -segment_format_options 
			movflags=+faststart -reset_timestamps 1 
			-vcodec h264 -q:v 4 -an -r 24 -strftime 1 %s `,
		},
	}

	gstreamer = processCmdTemplate{
		processor: "gst-launch-1.0",
		// tested with nano
		rtsp: cmdTemplate{
			base:    "-e --gst-debug-level=3 rtspsrc location=%s ! rtph264depay ! h264parse ! tee name=t ",
			capture: "t. ! queue ! avdec_h264 ! queue flush-on-eos=true ! videorate ! video/x-raw,framerate=%d/1 ! jpegenc ! multifilesink post-messages=true location=%s max-files=1 ",
			stream:  "t. ! queue ! flvmux streamable=true ! rtmpsink sync=false location=%s ",
			video:   "t. ! queue ! splitmuxsink max-size-time=%d location=%s ",
		},

		// tested with jx core board
		webcam: cmdTemplate{
			base:    "-e --gst-debug-level=3 v4l2src device=%s ",
			quality: "! image/jpeg,width=%[1]d,height=%[2]d,framerate=%[3]d/1 ! jpegdec ! tee name=t ",
			capture: "t. ! queue flush-on-eos=true ! videorate ! video/x-raw,framerate=%d/1 ! jpegenc ! multifilesink location=%s max-files=1 post-messages=true ",
			h264:    "t. ! queue ! videoconvert ! queue ! videoscale ! video/x-raw,width=%[1]d,height=%[2]d ! queue ! mpph264enc vbr=false bitrate=\"800000\" filerate=false ! queue ! h264parse ! tee name=v ",
			stream:  "v. ! queue ! flvmux streamable=true ! rtmpsink sync=false location=%s ",
			video:   "v. ! queue ! splitmuxsink max-size-time=%d location=%s ",
		},
	}
)
