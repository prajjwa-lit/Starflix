document.addEventListener("DOMContentLoaded", function () {
  const videoList = document.getElementById("video-list");
  const videoPlayer = document.getElementById("video-player");
  const player = document.getElementById("player");
  const videoTitle = document.getElementById("video-title");
  const backButton = document.getElementById("back-button");
  const videosTab = document.getElementById("videos-tab");
  const uploadTab = document.getElementById("upload-tab");
  const videosPage = document.getElementById("videos-page");
  const uploadPage = document.getElementById("upload-page");
  const uploadForm = document.getElementById("upload-form");
  const fileInput = document.getElementById("file-input");
  const selectedFile = document.getElementById("selected-file");
  const fileInfo = document.getElementById("file-info");
  const progressContainer = document.getElementById("progress-container");
  const progressBar = document.getElementById("progress-bar");
  const uploadMessage = document.getElementById("upload-message");
  const dropArea = document.getElementById("drop-area");
  videosTab.addEventListener("click", function () {
    videosTab.classList.add("active");
    uploadTab.classList.remove("active");
    videosPage.classList.add("active");
    uploadPage.classList.remove("active");
    loadVideos();
  });

  uploadTab.addEventListener("click", function () {
    uploadTab.classList.add("active");
    videosTab.classList.remove("active");
    uploadPage.classList.add("active");
    videosPage.classList.remove("active");
  });
  function formatFileSize(bytes) {
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    if (bytes === 0) return "0 B";
    const i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
    return Math.round(bytes / Math.pow(1024, i), 2) + " " + sizes[i];
  }
  function loadVideos() {
    videoList.innerHTML = '<div class="loading">Loading videos...</div>';

    fetch("/api/videos")
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to load videos");
        }
        return response.json();
      })
      .then((videos) => {
        videoList.innerHTML = "";

        if (videos.length === 0) {
          videoList.innerHTML = '<div class="loading">No videos found</div>';
          return;
        }

        videos.forEach((video) => {
          const videoItem = document.createElement("div");
          videoItem.className = "video-item";

          videoItem.innerHTML = `
                        <div class="video-thumbnail">▶</div>
                        <div class="video-info">
                            <div class="video-title">${video.name}</div>
                            <div class="video-size">${formatFileSize(video.size)}</div>
                        </div>
                    `;

          videoItem.addEventListener("click", () => {
            videoList.classList.add("hidden");
            videoPlayer.classList.remove("hidden");
            player.src = `/videos/${encodeURIComponent(video.path)}`;
            videoTitle.textContent = video.name;
            player
              .play()
              .catch((err) => console.error("Playback failed:", err));
          });

          videoList.appendChild(videoItem);
        });
      })
      .catch((error) => {
        console.error("Error:", error);
        videoList.innerHTML = `<div class="loading error">Error loading videos: ${error.message}</div>`;
      });
  }
  backButton.addEventListener("click", () => {
    player.pause();
    player.removeAttribute("src");
    player.load();
    videoPlayer.classList.add("hidden");
    videoList.classList.remove("hidden");
  });
  fileInput.addEventListener("change", () => {
    const file = fileInput.files[0];
    if (file) {
      fileInfo.classList.remove("hidden");
      selectedFile.textContent = `Selected: ${file.name} (${formatFileSize(file.size)})`;
    } else {
      fileInfo.classList.add("hidden");
    }
  });
  dropArea.addEventListener("dragover", (e) => {
    e.preventDefault();
    dropArea.classList.add("dragover");
  });

  dropArea.addEventListener("dragleave", () => {
    dropArea.classList.remove("dragover");
  });

  dropArea.addEventListener("drop", (e) => {
    e.preventDefault();
    dropArea.classList.remove("dragover");

    if (e.dataTransfer.files.length) {
      fileInput.files = e.dataTransfer.files;
      const file = fileInput.files[0];
      fileInfo.classList.remove("hidden");
      selectedFile.textContent = `Selected: ${file.name} (${formatFileSize(file.size)})`;
    }
  });

  dropArea.addEventListener("click", () => {
    fileInput.click();
  });
  uploadForm.addEventListener("submit", (e) => {
    e.preventDefault();

    const file = fileInput.files[0];
    if (!file) {
      uploadMessage.textContent = "Please select a file to upload";
      uploadMessage.className = "message error";
      return;
    }
    const videoTypes = [
      "video/mp4",
      "video/webm",
      "video/ogg",
      "video/quicktime",
      "video/x-msvideo",
    ];
    if (!videoTypes.includes(file.type)) {
      uploadMessage.textContent = "Please select a valid video file";
      uploadMessage.className = "message error";
      return;
    }

    const formData = new FormData();
    formData.append("file", file);
    uploadMessage.textContent = "";
    uploadMessage.className = "message";
    progressContainer.classList.remove("hidden");
    progressBar.style.width = "0%";
    const xhr = new XMLHttpRequest();

    xhr.open("POST", "/api/upload", true);

    xhr.upload.addEventListener("progress", (e) => {
      if (e.lengthComputable) {
        const percent = Math.round((e.loaded / e.total) * 100);
        progressBar.style.width = percent + "%";
      }
    });

    xhr.onload = function () {
      progressContainer.classList.add("hidden");

      if (xhr.status === 200) {
        const response = JSON.parse(xhr.responseText);
        uploadMessage.textContent = response.message;
        uploadMessage.className = "message success";
        uploadForm.reset();
        fileInfo.classList.add("hidden");
      } else {
        try {
          const response = JSON.parse(xhr.responseText);
          uploadMessage.textContent = response.message || "Upload failed";
        } catch (e) {
          uploadMessage.textContent = "Upload failed";
        }
        uploadMessage.className = "message error";
      }
    };

    xhr.onerror = function () {
      progressContainer.classList.add("hidden");
      uploadMessage.textContent = "Network error occurred";
      uploadMessage.className = "message error";
    };

    xhr.send(formData);
  });
  loadVideos();
});
