document.addEventListener("DOMContentLoaded", function () {
  const videosTab = document.getElementById("videos-tab");
  const uploadTab = document.getElementById("upload-tab");
  const videosPage = document.getElementById("videos-page");
  const uploadPage = document.getElementById("upload-page");
  const videoCategories = document.getElementById("video-categories");
  const videoPlayer = document.getElementById("video-player");
  const player = document.getElementById("player");
  const videoTitle = document.getElementById("video-title");
  const videoDetails = document.getElementById("video-details");
  const backButton = document.getElementById("back-button");
  const loadingIndicator = document.getElementById("loading-indicator");
  const featuredSection = document.getElementById("featured-section");
  const featuredTitle = document.getElementById("featured-title");
  const featuredDescription = document.getElementById("featured-description");
  const featuredPlay = document.getElementById("featured-play");
  const featuredInfo = document.getElementById("featured-info");
  const genreNav = document.getElementById("genre-nav");
  const uploadForm = document.getElementById("upload-form");
  const fileInput = document.getElementById("file-input");
  const coverInput = document.getElementById("cover-input");
  const selectedFile = document.getElementById("selected-file");
  const selectedCover = document.getElementById("selected-cover");
  const fileInfo = document.getElementById("file-info");
  const progressContainer = document.getElementById("progress-container");
  const progressBar = document.getElementById("progress-bar");
  const uploadMessage = document.getElementById("upload-message");
  const dropArea = document.getElementById("drop-area");
  const coverUploadArea = document.getElementById("cover-upload-area");
  const genreSelect = document.getElementById("genre");
  const movieModal = document.getElementById("movie-modal");
  const modalTitle = document.getElementById("modal-title");
  const modalPoster = document.getElementById("modal-poster");
  const modalDescription = document.getElementById("modal-description");
  const modalGenre = document.getElementById("modal-genre");
  const modalYear = document.getElementById("modal-year");
  const modalSize = document.getElementById("modal-size");
  const modalPlay = document.getElementById("modal-play");
  const modalClose = document.getElementById("modal-close");
  let currentGenre = "all";
  let allVideos = [];
  let featuredVideo = null;
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
    loadGenres();
  });
  function formatFileSize(bytes) {
    if (bytes === 0) return "0 B";
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return Math.round(bytes / Math.pow(1024, i), 2) + " " + sizes[i];
  }
  function getReleaseYear(video) {
    return video.release_year > 0 ? video.release_year : "";
  }
  function getCoverImageUrl(video) {
    if (video.cover_image && video.cover_image !== "") {
      return `/covers/${encodeURIComponent(video.cover_image)}`;
    }
    const defaultCovers = {
      Action: "action-default.jpg",
      Comedy: "comedy-default.jpg",
      Drama: "drama-default.jpg",
      Documentary: "documentary-default.jpg",
      Thriller: "thriller-default.jpg",
      Horror: "horror-default.jpg",
      "Sci-Fi": "scifi-default.jpg",
    };

    if (video.genre && defaultCovers[video.genre]) {
      return `/covers/defaults/${defaultCovers[video.genre]}`;
    }
    const colors = [
      "linear-gradient(45deg, #E50914, #B20710)",
      "linear-gradient(45deg, #0F79AF, #0C5A83)",
      "linear-gradient(45deg, #8C4AFF, #6D37C8)",
      "linear-gradient(45deg, #FF6B00, #C65200)",
      "linear-gradient(45deg, #00AD85, #008F6D)",
    ];
    const hash = video.title
      .split("")
      .reduce((acc, char) => acc + char.charCodeAt(0), 0);
    const colorIndex = hash % colors.length;
    return `data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='100' height='100' viewBox='0 0 100 100'%3E%3Crect width='100' height='100' fill='${encodeURIComponent(colors[colorIndex])}'/%3E%3C/svg%3E`;
  }
  function loadGenres() {
    fetch("/api/genres")
      .then((response) => {
        if (!response.ok) throw new Error("Failed to load genres");
        return response.json();
      })
      .then((genres) => {
        const genreNavContent = `
          <button class="genre-btn active" data-genre="all">All</button>
          ${genres.map((g) => `<button class="genre-btn" data-genre="${g.name}">${g.name}</button>`).join("")}
        `;
        genreNav.innerHTML = genreNavContent;
        document.querySelectorAll(".genre-btn").forEach((btn) => {
          btn.addEventListener("click", () => {
            document
              .querySelectorAll(".genre-btn")
              .forEach((b) => b.classList.remove("active"));
            btn.classList.add("active");
            currentGenre = btn.dataset.genre;
            filterVideosByGenre(currentGenre);
          });
        });
        let options = '<option value="">Select genre</option>';
        genres.forEach((g) => {
          options += `<option value="${g.name}">${g.name}</option>`;
        });
        genreSelect.innerHTML = options;
      })
      .catch((error) => {
        console.error("Error loading genres:", error);
      });
  }
  function loadVideos() {
    loadingIndicator.classList.remove("hidden");
    videoCategories.innerHTML = "";

    fetch("/api/videos")
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to load videos");
        }
        return response.json();
      })
      .then((videos) => {
        loadingIndicator.classList.add("hidden");
        allVideos = videos;

        if (videos.length === 0) {
          videoCategories.innerHTML =
            '<div class="loading">No videos found</div>';
          return;
        }
        selectFeaturedVideo(videos);
        organizeVideosByGenre(videos);
      })
      .catch((error) => {
        console.error("Error:", error);
        loadingIndicator.classList.add("hidden");
        videoCategories.innerHTML = `<div class="loading error">Error loading videos: ${error.message}</div>`;
      });
  }
  function selectFeaturedVideo(videos) {
    if (videos.length === 0) return;
    featuredVideo =
      videos.find((v) => v.cover_image && v.cover_image !== "") || videos[0];
    featuredSection.classList.remove("hidden");
    featuredTitle.textContent = featuredVideo.title;
    featuredDescription.textContent =
      featuredVideo.description ||
      `${featuredVideo.filename} - ${formatFileSize(featuredVideo.size)}`;
    const coverUrl = getCoverImageUrl(featuredVideo);
    featuredSection.style.backgroundImage = `url('${coverUrl}')`;
    featuredPlay.addEventListener("click", () => playVideo(featuredVideo));
    featuredInfo.addEventListener("click", () =>
      showVideoDetails(featuredVideo),
    );
  }
  function organizeVideosByGenre(videos) {
    const genres = [
      ...new Set(videos.filter((v) => v.genre).map((v) => v.genre)),
    ];
    let categoriesHTML = `
      <div class="category" id="category-all">
        <h2 class="category-header">All Videos</h2>
        <div class="video-slider">


        ${renderVideoItems(videos)}
        </div>
      </div>
    `;
    genres.forEach((genre) => {
      const genreVideos = videos.filter((v) => v.genre === genre);
      if (genreVideos.length > 0) {
        categoriesHTML += `
          <div class="category" id="category-${genre.toLowerCase().replace(/\s+/g, "-")}">
            <h2 class="category-header">${genre}</h2>
            <div class="video-slider">
              ${renderVideoItems(genreVideos)}
            </div>
          </div>
        `;
      }
    });
    const ungenredVideos = videos.filter((v) => !v.genre);
    if (ungenredVideos.length > 0) {
      categoriesHTML += `
        <div class="category" id="category-uncategorized">
          <h2 class="category-header">Uncategorized</h2>
          <div class="video-slider">
            ${renderVideoItems(ungenredVideos)}
          </div>
        </div>
      `;
    }

    videoCategories.innerHTML = categoriesHTML;
    document.querySelectorAll(".video-item").forEach((item) => {
      const videoId = parseInt(item.dataset.id);
      const video = allVideos.find((v) => v.id === videoId);

      item.addEventListener("click", () => showVideoDetails(video));
    });
  }
  function filterVideosByGenre(genre) {
    if (genre === "all") {
      document.querySelectorAll(".category").forEach((cat) => {
        cat.classList.remove("hidden");
      });
    } else {
      document.querySelectorAll(".category").forEach((cat) => {
        if (
          cat.id === `category-${genre.toLowerCase().replace(/\s+/g, "-")}` ||
          cat.id === "category-all"
        ) {
          cat.classList.remove("hidden");
        } else {
          cat.classList.add("hidden");
        }
      });
    }
  }

  function renderVideoItems(videos) {
    return videos
      .map((video) => {
        const coverUrl = getCoverImageUrl(video);
        return `
        <div class="video-item" data-id="${video.id}">
          <div class="video-thumbnail" style="background-image: url('${coverUrl}')">
            <div class="play-indicator">
              <svg viewBox="0 0 24 24">
                <path d="M8 5v14l11-7z"></path>
              </svg>
            </div>
          </div>
          <div class="video-info">
            <div class="video-title">${video.title}</div>
            <div class="video-meta">
              <span class="video-genre">${video.genre || ""}</span>
              <span class="video-year">${getReleaseYear(video)}</span>
            </div>
            <div class="video-description">${video.description || ""}</div>
          </div>
        </div>
      `;
      })
      .join("");
  }
  function playVideo(video) {
    videoPlayer.classList.remove("hidden");
    featuredSection.classList.add("hidden");
    videoCategories.classList.add("hidden");
    genreNav.classList.add("hidden");
    player.src = `/videos/${encodeURIComponent(video.path)}`;
    videoTitle.textContent = video.title;
    videoDetails.innerHTML = `
      <div class="video-meta-details">
        <div class="detail-item">
          <span class="detail-label">Genre:</span>
          <span class="detail-value">${video.genre || "Not specified"}</span>
        </div>
        ${
          video.release_year
            ? `
          <div class="detail-item">
            <span class="detail-label">Year:</span>
            <span class="detail-value">${video.release_year}</span>
          </div>
        `
            : ""
        }
        <div class="detail-item">
          <span class="detail-label">Size:</span>
          <span class="detail-value">${formatFileSize(video.size)}</span>
        </div>
      </div>
      ${video.description ? `<p class="video-description">${video.description}</p>` : ""}
    `;
    player.play().catch((err) => console.error("Playback failed:", err));
    window.scrollTo(0, 0);
  }
  function showVideoDetails(video) {
    modalTitle.textContent = video.title;
    modalDescription.textContent =
      video.description || "No description available.";
    modalGenre.textContent = video.genre || "Not specified";
    modalYear.textContent =
      video.release_year > 0 ? video.release_year : "Not specified";
    modalSize.textContent = formatFileSize(video.size);
    const coverUrl = getCoverImageUrl(video);
    modalPoster.src = coverUrl;
    modalPlay.onclick = () => {
      closeModal();
      playVideo(video);
    };
    movieModal.style.display = "block";
    document.body.style.overflow = "hidden";
  }
  function closeModal() {
    movieModal.style.display = "none";
    document.body.style.overflow = "auto";
  }
  modalClose.addEventListener("click", closeModal);
  window.addEventListener("click", (e) => {
    if (e.target === movieModal) {
      closeModal();
    }
  });
  backButton.addEventListener("click", () => {
    player.pause();
    player.removeAttribute("src");
    player.load();
    videoPlayer.classList.add("hidden");
    featuredSection.classList.remove("hidden");
    videoCategories.classList.remove("hidden");
    genreNav.classList.remove("hidden");
  });
  fileInput.addEventListener("change", () => {
    const file = fileInput.files[0];
    if (file) {
      fileInfo.classList.remove("hidden");
      selectedFile.textContent = `Selected: ${file.name} (${formatFileSize(file.size)})`;
      const titleInput = document.getElementById("title");
      if (!titleInput.value) {
        const filename = file.name.replace(/\.[^/.]+$/, "");
        titleInput.value = filename;
      }
    } else {
      fileInfo.classList.add("hidden");
    }
  });
  coverInput.addEventListener("change", () => {
    const file = coverInput.files[0];
    if (file) {
      selectedCover.classList.remove("hidden");
      selectedCover.textContent = `Selected: ${file.name}`;
      const reader = new FileReader();
      reader.onload = (e) => {
        coverUploadArea.style.backgroundImage = `url('${e.target.result}')`;
        coverUploadArea.style.backgroundSize = "cover";
        coverUploadArea.style.backgroundPosition = "center";
        coverUploadArea.querySelector("svg").style.display = "none";
        coverUploadArea.querySelector("p").style.display = "none";
      };
      reader.readAsDataURL(file);
    } else {
      selectedCover.classList.add("hidden");
      coverUploadArea.style.backgroundImage = "";
      coverUploadArea.querySelector("svg").style.display = "block";
      coverUploadArea.querySelector("p").style.display = "block";
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
      const titleInput = document.getElementById("title");
      if (!titleInput.value) {
        const filename = file.name.replace(/\.[^/.]+$/, "");
        titleInput.value = filename;
      }
    }
  });
  dropArea.addEventListener("click", () => {
    fileInput.click();
  });
  coverUploadArea.addEventListener("click", () => {
    coverInput.click();
  });
  uploadForm.addEventListener("submit", (e) => {
    e.preventDefault();

    const file = fileInput.files[0];
    if (!file) {
      uploadMessage.textContent = "Please select a video file to upload";
      uploadMessage.className = "message error";
      return;
    }
    const videoTypes = [
      "video/mp4",
      "video/webm",
      "video/ogg",
      "video/quicktime",
      "video/x-msvideo",
      "video/mp2t",
    ];
    if (
      !videoTypes.includes(file.type) &&
      !file.name.toLowerCase().endsWith(".ts")
    ) {
      uploadMessage.textContent = "Please select a valid video file";
      uploadMessage.className = "message error";
      return;
    }

    const formData = new FormData(uploadForm);
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
        let errorMessage = "Upload failed";
        try {
          const response = JSON.parse(xhr.responseText);
          errorMessage = `Upload failed: ${response.message || xhr.statusText}`;
        } catch (e) {
          errorMessage = `Upload failed: ${xhr.statusText}`;
        }
        console.error("Upload error:", errorMessage);
        uploadMessage.textContent = errorMessage;
        uploadMessage.className = "message error";
      }
    };

    xhr.onerror = function (e) {
      console.error("Network error:", e);
      progressContainer.classList.add("hidden");
      uploadMessage.textContent = `Network error occurred: ${e.message || "Unknown error"}`;
      uploadMessage.className = "message error";
    };
    xhr.send(formData);
  });
  document.addEventListener("keydown", function (e) {
    if (!videoPlayer.classList.contains("hidden")) {
      if (e.code === "Space") {
        e.preventDefault();
        if (player.paused) {
          player.play();
        } else {
          player.pause();
        }
      }
      if (e.code === "Escape") {
        backButton.click();
      }
    }
    if (e.code === "Escape" && movieModal.style.display === "block") {
      closeModal();
    }
  });
  const header = document.querySelector("header");
  window.addEventListener("scroll", () => {
    if (window.scrollY > 50) {
      header.classList.add("scrolled");
    } else {
      header.classList.remove("scrolled");
    }
  });
  loadVideos();
});
