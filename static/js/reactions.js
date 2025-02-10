import { showNotification } from "./components/notifications.js";

export function attachInteractionListeners() {
  const reactionsBtns = document.querySelectorAll(".like-btn, .dislike-btn");

  reactionsBtns.forEach((button) => {
    if (button.classList.contains("comment")) {
      console.log(
        `Attaching event listener to comment button with ID: ${button.dataset.commentId}`
      );
      button.addEventListener("click", (e) => handleInteraction(e, "comment"));
    } else if (button.classList.contains("post")) {
      console.log(
        `Attaching event listener to post button with ID: ${button.dataset.postId}`
      );
      button.addEventListener("click", (e) => handleInteraction(e, "post"));
    } else {
      console.warn("Invalid button type detected:", button);
      showNotification("Invalid button type", "error");
    }
  });
}

export async function handleInteraction(e, itemtype) {
  e.preventDefault();
  const button = e.target;
  const action = button.value?.toLowerCase();

  console.log(`Interaction detected: ${action} on a ${itemtype}`);

  let item_id = null;
  if (itemtype === "post") {
    item_id = button.dataset.postId;
  } else if (itemtype === "comment") {
    item_id = button.dataset.commentId;
  } else {
    console.error("Invalid item type:", itemtype);
    showNotification("Invalid button type", "error");
    return;
  }

  if (action != "like" && action != "dislike") {
    console.error("Invalid action:", action);
    showNotification("Invalid action", "error");
    return;
  }

  try {
    console.log(
      `Sending request to /api/reactions with item_id: ${item_id}, item_type: ${itemtype}, action: ${action}`
    );

    const response = await fetch("/api/reactions", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        item_id: item_id,
        item_type: itemtype,
        action: action,
      }),
      credentials: "include",
    });

    if (response.ok) {
      const result = await response.json();                                

      if (itemtype === "post") {
        document.querySelector(
          `#like-postcount-${result.item_id}`
        ).textContent = `${result.likes} Likes`;
        document.querySelector(
          `#dislike-postcount-${result.item_id}`
        ).textContent = `${result.dislikes} Dislikes`;
      } else if (itemtype === "comment") {
        document.querySelector(
          `#likes-commentcount-${result.item_id}`
        ).textContent = `${result.likes} Likes`;
        document.querySelector(
          `#dislikes-commentcount-${result.item_id}`
        ).textContent = `${result.dislikes} Dislikes`;
      } else {
        showNotification(`Unkown Type: Reactions only on Posts and Comments`, "error");
        return
      }
     
    } else {
      console.error("Failed to update interaction. Status:", response.status);
      const err = await response.json();
      console.error("Failed to update interaction. Status:", response.status);
      showNotification(err.message, "error");
    }
  } catch (error) {
    console.error("Error updating interaction:", error);
    showNotification("Error updating interaction", "error");
  }
}
