const local_vid_div = document.getElementById('local-video');
const remote_vids_div = document.getElementById('remote-videos');
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

    switch (device) {
        case 'camera':
            const video = document.createElement('video')
            try {
                const media = await IonSDK.LocalStream.getUserMedia({
                    resolution: 'vga',
                    audio: true,
                    codec: 'vp8',
    
                })

                video.autoplay = true;
                video.id = 'user-video'
                video.srcObject = media; 

                local_vid_div.insertAdjacentElement('afterbegin', video);

                clientLocal.publish(media);
                mediaControls.classList.remove('hidden');
            }catch(error) {
                console.log('camera error:', error)
            }
            break;
        case 'screen':
            const screen = document.createElement('video')
            try {
                const media = await IonSDK.LocalStream.getDisplayMedia({
                    resolution: 'vga',
                    audio: true,
                    codec: 'vp8',
    
                })

                screen.autoplay = true;
                screen.srcObject = media; 

                local_vid_div.insertAdjacentElement('afterbegin', screen);

                clientLocal.publish(media);
            }catch(error) {
                console.log('screen error:', error)
            }
    }
}

show('camera');

clientLocal.ontrack = (track, stream) => {
    const video = document.createElement('video');
    if (track.kind === 'video') {
        track.onunmute = () => { // if there is a track
            video.id = track.id;
            video.srcObject = stream;
            video.autoplay = true;
            video.controls = false;

            console.log('adding track:', track.id)
            remote_vids_div.appendChild(video)

            stream.onremovetrack = (e) => {
                console.log('about to remove:', e.track.kind);
                if (e.track.kind === 'video') {
                    const videoToRemove = document.getElementById(e.track.id);
                    remote_vids_div.removeChild(videoToRemove);
                }
            }
        }
    }
    
}

const toggleMicrophone = () => {
    const audio = document.getElementById('user-video').srcObject.getTrack('audio');
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
    const video = document.getElementById('user-video').srcObject.getTrack('video');
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