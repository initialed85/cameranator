import time
import traceback
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
        super().__init__(
            camera=camera,
            input_path=input_path,
            output_path=output_path,
            output_fps=output_fps,
            label=label,
            output_fourcc=output_fourcc,
            output_extension=output_extension,
        )

    def __iter__(self):
        start = time.time()

        process_fpses = []

        while True:
            self.frame_counter += 1

            ret, frame = self.video_capture.read()
            if ret is False or frame is None:
                break

            process_fps = self.frame_counter / (time.time() - start)
            process_fpses.append(process_fps)

            yield frame

        min_process_fps = 0
        avg_process_fps = 0
        max_process_fps = 0

        if len(process_fpses) > 0:
            min_process_fps = min(process_fpses)
            avg_process_fps = sum(process_fpses) / len(process_fpses)
            max_process_fps = max(process_fpses)

        print(
            f"fps_min={min_process_fps}, fps_avg={avg_process_fps}, fps_max={max_process_fps}"
        )

        if self.output_video is not None:
            self.output_video.release()
            print(f"Output video file saved to: {self.get_output_file_path()}")

        self.video_capture.release()
        cv2.destroyAllWindows()

    def _write(self, frame: np.ndarray):
        try:
            if self.output_video is None:
                output_file_path = self.get_output_file_path()

                fourcc = cv2.VideoWriter_fourcc(
                    *self.get_codec_fourcc(output_file_path)
                )

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
        except Exception:
            traceback.print_exc()
            raise

    def write(self, frame: np.ndarray) -> int:
        return -1
