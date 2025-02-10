import { showNotification } from "./components/notifications.js";
import { handleInteraction } from "./reactions.js";

let lastClickedTime = 0;
const THROTTLE_DELAY = 1000;
function createCommentMarkup(comment) {
  const commentItem = document.createElement("li");
  const commentDiv = document.createElement("div");
  commentDiv.classList.add("comment");

  const commentText = document.createElement("p");
  const strongElement = document.createElement("strong");
  strongElement.textContent = comment.Username || "Anonymous";
  const commentContent = document.createTextNode(`: ${comment.Content}`);

  commentText.appendChild(strongElement);
  commentText.appendChild(commentContent);

  const likesCount = document.createElement("button");
  likesCount.classList.add("likes-count", "like-btn");
  likesCount.textContent = `${comment.Likes || 0} Likes`;
  likesCount.id = `likes-commentcount-${comment.ID}`;

  const dislikesCount = document.createElement("button");
  dislikesCount.classList.add("dislikes-count", "dislike-btn");
  dislikesCount.textContent = `${comment.Dislikes || 0} Dislikes`;
  dislikesCount.id = `dislikes-commentcount-${comment.ID}`;

  likesCount.dataset.commentId = comment.ID;
  likesCount.value = "like";
  likesCount.addEventListener("click", (e) => {
    const now = Date.now();
    if (now - lastClickedTime >= THROTTLE_DELAY) {
      lastClickedTime = now;
      handleInteraction(e, "comment");
    }
  });

  dislikesCount.dataset.commentId = comment.ID;
  dislikesCount.value = "dislike";
  dislikesCount.addEventListener("click", (e) => {
    const now = Date.now();
    if (now - lastClickedTime >= THROTTLE_DELAY) {
      lastClickedTime = now;
      handleInteraction(e, "comment");
    }
  });

  const interactionsContainer = document.createElement("div");
  interactionsContainer.classList.add("interactions-container");
  interactionsContainer.appendChild(likesCount);
  interactionsContainer.appendChild(dislikesCount);
  commentDiv.appendChild(commentText);
  commentDiv.appendChild(interactionsContainer);

  commentItem.appendChild(commentDiv);
  return commentItem;
}

export function CommentHandler(
  commentButton,
  commentTextArea,
  commentList,
  postId
) {
  commentButton.addEventListener("click", (e) => {
    e.preventDefault();
    AddComment(commentTextArea, commentList, postId);
  });
}

export async function AddComment(commentTextArea, commentList, postId) {
  const commentContent = commentTextArea.value.trim();

  if (!commentContent) {
    showNotification("Comment cannot be empty", "error");
    return;
  }

  try {
    const response = await fetch(`/api/comments/${postId}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ content: commentContent, post_id: postId }),
    });

    if (response.status === 201) {
      const newComment = await response.json();
      const newCommentElement = createCommentElement(newComment);
      commentList.appendChild(newCommentElement);
      commentTextArea.value = "";
      showNotification("Comment added successfully", "success");
    } else {
      const errorMsg = await response.json();
      showNotification(errorMsg.message, "error");
    }
  } catch (error) {
    showNotification("Error submitting comment", "error");
  }
}

export function createCommentElement(comment) {
  return createCommentMarkup(comment);
}
