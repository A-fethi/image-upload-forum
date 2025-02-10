import { showNotification } from "./components/notifications.js";
import { createPostElement } from "./posts.js";

export let post = {};

export function checkPost() {
  const isLoggedIn = document.cookie.includes("IsLoggedIn=true");
  const createPostButton = document.getElementById("create-post-button");

  if (createPostButton) {
    createPostButton.addEventListener("click", async (event) => {
      event.preventDefault();
      const titleInput = document.querySelector('input[name="title"]');
      const contentInput = document.querySelector('textarea[name="content"]');
      const imageInput = document.querySelector('input[name="image"]');

      const selectedCategories = Array.from(
        document.querySelectorAll('input[name="category"]:checked')
      ).map((checkbox) => checkbox.value);

      if (titleInput && contentInput) {
        let title = titleInput.value.trim();
        let content = contentInput.value.trim();
        let imageFile = imageInput.files[0]; // Get the actual file

        if (!title || !content) {
          showNotification("Error: Title and Content cannot be empty", "error");
          return;
        }

        // Create FormData object to handle file upload
        const formData = new FormData();
        formData.append("title", title);
        formData.append("content", content);
        if (imageFile) {
          formData.append("image", imageFile); // Only append if an image is uploaded
        }
        selectedCategories.forEach((category) =>
          formData.append("categories", category)
        );

        try {
          const resp = await fetch("/api/posts/add", {
            method: "POST",
            body: formData,
            credentials: "include",
          });

          if (resp.status === 201) {
            const responseData = await resp.json();
            titleInput.value = "";
            contentInput.value = "";
            imageInput.value = ""; // Reset file input

            document
              .querySelectorAll('input[name="category"]:checked')
              .forEach((checkbox) => (checkbox.checked = false));

            const postsElement = document.getElementById("posts-container");
            postsElement.prepend(createPostElement(responseData));
            showNotification("Post created successfully!", "success");
          } else {
            const responseData = await resp.json();
            console.error("Failed to create post:", resp.statusText);
            showNotification(responseData.message, "error");
          }
        } catch (error) {
          console.error("Error occurred while creating post:", error);
          showNotification("An error occurred, Please try again later", "error");
        }
      } else {
        console.error("Title or Content inputs not found.");
        showNotification("Error: Title/Content cannot be empty", "error");
      }
    });
  } else {
    console.error("Submit button not found.");
  }
}
