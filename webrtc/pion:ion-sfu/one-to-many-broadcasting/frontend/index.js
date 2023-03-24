const pub_video = document.getElementById('pub_video');
const sub_video = document.getElementById('sub_video');
const cam_btn = document.getElementById('btn_cam');
const screen_btn = document.getElementById('btn_screen');
const mediaControls = document.getElementById('media-controls');

const serverURL = 'ws://localhost:7000/ws'

const config = {
    iceServers: [
        {urls: "stun:stun.l.google.com:19302"},
    ]
}

const signalLocal = new Signal.IonSFUJSONRPCSignal(serverURL);
const clientLocal = new IonSDK.Client(signalLocal, config);

signalLocal.onopen = () => clientLocal.join("some-id");

const show = async (device) => {
    if (pub_video && pub_video.srcObject !== null) {
        pub_video.srcObject = null;
    }
    
    if (sub_video && sub_video.srcObject !== null) {
        sub_video.srcObject = null;
    }

    switch (device) {
        case 'camera':
            try {
                const media = await IonSDK.LocalStream.getUserMedia({
                    resolution: 'vga',
                    audio: true,
                    codec: 'vp8',
    
                })

                pub_video.srcObject = media; 
                pub_video.controls = false;
                // cam_btn.disabled = true;
                // screen_btn.disabled = true;

                clientLocal.publish(media);

                mediaControls.classList.remove('hidden');
            }catch(error) {
                console.log('camera error:', error)
            }
            break;
        case 'screen':
            try {
                const media = await IonSDK.LocalStream.getDisplayMedia({
                    resolution: 'vga',
                    audio: true,
                    codec: 'vp8',
    
                })

                pub_video.srcObject = media; 
                pub_video.controls = false;
                // cam_btn.disabled = true;
                // screen_btn.disabled = true;

                clientLocal.publish(media);

                mediaControls.classList.remove('hidden');
            }catch(error) {
                console.log('screen error:', error)
            }
    }
}

show('camera');

clientLocal.ontrack = (track, stream) => {
    track.onunmute = () => { // if there is a track
        sub_video.srcObject = stream;
        sub_video.autoplay = true;
        sub_video.controls = false;

        stream.onremovetrack = () => {
            sub_video.srcObject = null;
        }
    }
}

const toggleMicrophone = () => {
    const audio = pub_video.srcObject.getTrack('audio');
    if (audio.enabled) {
        audio.enabled = false;
        document.getElementById('mic-on-icon').classList.add('hidden')
        document.getElementById('mic-off-icon').classList.remove('hidden')
    }else {
        audio.enabled = true;
        document.getElementById('mic-off-icon').classList.add('hidden')
        document.getElementById('mic-on-icon').classList.remove('hidden')
    }
}

const toggleCamera = () => {
    const video = pub_video.srcObject.getTrack('video');
    if (video.enabled) {
        video.enabled = false;
        document.getElementById('cam-on-icon').classList.add('hidden')
        document.getElementById('cam-off-icon').classList.remove('hidden')
    }else {
        video.enabled = true;
        document.getElementById('cam-off-icon').classList.add('hidden')
        document.getElementById('cam-on-icon').classList.remove('hidden')
    }
}