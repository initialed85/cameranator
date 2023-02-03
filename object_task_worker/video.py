# credit to https://github.com/tryolabs/norfair/blob/master/norfair/video.py; extended w/ ThreadPoolExcutor

import time
from concurrent.futures import ThreadPoolExecutor, wait
from typing import Optional

import numpy as np
from norfair import Video as _Video
from norfair.utils import DummyOpenCVImport

try:
    import cv2
except ImportError:
    cv2 = DummyOpenCVImport()


class Video(_Video):
    def __init__(
        self,
        camera: Optional[int] = None,
        input_path: Optional[str] = None,
        output_path: str = ".",
        output_fps: Optional[float] = None,
        label: str = "",
        output_fourcc: Optional[str] = None,
        output_extension: str = "mp4",
    ):
        self._executor = ThreadPoolExecutor(max_workers=1)
        self._futures = []

        super().__init__(
            camera=camera,
            input_path=input_path,
            output_path=output_path,
            output_fps=output_fps,
            label=label,
            output_fourcc=output_fourcc,
            output_extension=output_extension,
        )

    def __del__(self):
        self._executor.shutdown(wait=True)
        del self

    def __iter__(self):
        with self.progress_bar as progress_bar:
            start = time.time()

            while True:
                self.frame_counter += 1
                ret, frame = self.video_capture.read()
                if ret is False or frame is None:
                    break
                process_fps = self.frame_counter / (time.time() - start)
                progress_bar.update(
                    self.task, advance=1, refresh=True, process_fps=process_fps
                )
                yield frame

        wait(self._futures)

        if self.output_video is not None:
            self.output_video.release()
            print(
                f"[white]Output video file saved to: {self.get_output_file_path()}[/white]"
            )
        self.video_capture.release()
        cv2.destroyAllWindows()

    def _write(self, frame: np.ndarray):
        if self.output_video is None:
            output_file_path = self.get_output_file_path()
            fourcc = cv2.VideoWriter_fourcc(*self.get_codec_fourcc(output_file_path))
            output_size = (
                frame.shape[1],
                frame.shape[0],
            )
            self.output_video = cv2.VideoWriter(
                output_file_path,
                fourcc,
                self.output_fps,
                output_size,
            )

        self.output_video.write(frame)

    def write(self, frame: np.ndarray) -> int:
        self._futures.append(self._executor.submit(self._write, frame))
        return -1
