import { createCommentElement } from "./comments.js"
import { handleInteraction } from "./reactions.js";
import { CommentHandler } from "./comments.js";

export function createPostElement(post) {
    const postElement = document.createElement('div');
    postElement.classList.add('post');

    const postHeader = document.createElement('div');
    postHeader.classList.add('post-header');

    const postAuthor = document.createElement('div');
    postAuthor.classList.add('post-author');

    const authorImg = document.createElement('img');
    authorImg.src = "https://www.w3schools.com/w3images/avatar2.png";
    authorImg.alt = "Author";
    authorImg.classList.add('author-avatar');

    const authorInfo = document.createElement('div');
    authorInfo.classList.add('author-info');

    const authorName = document.createElement('span');
    authorName.classList.add('author-name');
    authorName.textContent = post.Username;

    const postDate = document.createElement('span');
    postDate.classList.add('post-date');
    postDate.textContent = `• ${post.CreatedAt}`;

    const postImage = document.createElement('div')
    postImage.classList.add('post-image')
    
    if (post.ImageContent) {
        const postImgElement = document.createElement('img');
        postImgElement.src = `/static/uploads/${post.ImageContent}`; // Adjust based on your server
        postImgElement.alt = "Post Image";
        postImgElement.classList.add('post-img');
        postImage.appendChild(postImgElement);
    }
    
    console.log("imageeee",post.ImageContent);
    
    authorInfo.appendChild(authorName);
    authorInfo.appendChild(postDate);
    
    postAuthor.appendChild(authorImg);
    postAuthor.appendChild(authorInfo);
    postHeader.appendChild(postAuthor);
    
    const postTitle = document.createElement('h2');
    postTitle.textContent = post.Title;
    
    postHeader.appendChild(postTitle);
    postElement.appendChild(postHeader);
    
    const postContent = document.createElement('div');
    postContent.classList.add('post-content');
    postContent.textContent = post.Content;
    postElement.appendChild(postContent);
    postElement.appendChild(postImage);

    const categoriesSpan = document.createElement('span');
    categoriesSpan.classList.add('post-categories');
    categoriesSpan.dataset.categories = post.Categories;
    postElement.appendChild(categoriesSpan);
    if(post.Categories) {
        post.Categories?.split(" ").map((category) => {
            const span = document.createElement("span");
            span.textContent = category;
            span.classList.add(category);
            categoriesSpan.appendChild(span);
        });
    }
    
    const interactionsDiv = document.createElement('div');
    interactionsDiv.classList.add('post-interactions');

    const likeButton = document.createElement('button');
    likeButton.classList.add('interaction-button', 'like-btn');
    likeButton.dataset.postId = post.ID;
    likeButton.value = "like";
    likeButton.addEventListener("click", (e) => handleInteraction(e, "post"));

    const likeIcon = document.createElement('i');
    likeIcon.classList.add('fas', 'fa-thumbs-up');
    likeIcon.value = "like";
    likeIcon.dataset.postId = post.ID;

    const likeText = document.createElement('span');
    likeText.value = "like";
    likeText.dataset.postId = post.ID;
    likeText.id = `like-postcount-${post.ID}`;
    likeText.textContent = `${post.Likes} Likes`;

    likeButton.appendChild(likeIcon);
    likeButton.appendChild(likeText);
    interactionsDiv.appendChild(likeButton);

    const dislikeButton = document.createElement('button');
    dislikeButton.classList.add('interaction-button', 'dislike-btn');
    dislikeButton.dataset.postId = post.ID;
    dislikeButton.value = "dislike";
    dislikeButton.addEventListener("click", (e) => handleInteraction(e, "post"));

    const dislikeIcon = document.createElement('i');
    dislikeIcon.classList.add('fas', 'fa-thumbs-down');
    dislikeIcon.dataset.postId = post.ID;
    dislikeIcon.value = "dislike";

    const dislikeText = document.createElement('span');
    dislikeText.id = `dislike-postcount-${post.ID}`;
    dislikeText.dataset.postId = post.ID;
    dislikeText.value = "dislike";
    dislikeText.textContent = `${post.Dislikes} Dislikes`;

    dislikeButton.appendChild(dislikeIcon);
    dislikeButton.appendChild(dislikeText);
    interactionsDiv.appendChild(dislikeButton);

    const toggleCommentsButton = document.createElement('button');
    toggleCommentsButton.classList.add('interaction-button');
    toggleCommentsButton.classList.add('XYZ-Comments');
    toggleCommentsButton.dataset.postId = post.ID;

    const toggleCommentsIcon = document.createElement('i');
    toggleCommentsIcon.classList.add('fas', 'fa-comments');
    toggleCommentsIcon.classList.add('XYZ-Comments');
    toggleCommentsIcon.dataset.postId = post.ID;

    const toggleCommentsText = document.createElement('span');
    toggleCommentsText.textContent = "Show Comments";
    toggleCommentsText.classList.add('toggle-comments-text');
    toggleCommentsText.classList.add('XYZ-Comments');
    toggleCommentsText.dataset.postId = post.ID;

    toggleCommentsButton.appendChild(toggleCommentsIcon);
    toggleCommentsButton.appendChild(toggleCommentsText);
    interactionsDiv.appendChild(toggleCommentsButton);
    const commentsSection = document.createElement('div');
    commentsSection.classList.add('comments-section');
    commentsSection.style.display = "none";

    const commentsHeader = document.createElement('div');
    commentsHeader.classList.add('comments-header');

    const commentsIcon = document.createElement('i');
    commentsIcon.classList.add('fas', 'fa-comments');

    const commentsText = document.createElement('span');
    commentsText.textContent = `Comments`;

    commentsHeader.appendChild(commentsIcon);
    commentsHeader.appendChild(commentsText);

    commentsSection.appendChild(commentsHeader);

    const commentForm = document.createElement('div');
    commentForm.classList.add('comment-form');

    const commentTextArea = document.createElement('textarea');
    commentTextArea.name = "comment";
    commentTextArea.classList.add('comment-input');
    commentTextArea.placeholder = "Write a comment...";
    commentTextArea.required = true;

    const commentButton = document.createElement('button');
    commentButton.classList.add('comment-button');
    commentButton.textContent = "Comment";
    commentButton.dataset.postId = post.ID;

    commentForm.appendChild(commentTextArea);
    commentForm.appendChild(commentButton);
    commentsSection.appendChild(commentForm);

    const commentList = document.createElement('div');
    commentList.classList.add('comment-list');
    commentsSection.appendChild(commentList);
    commentsSection.dataset.postId = post.ID;
    post.Comments?.forEach(comment => {
        const commentElement = createCommentElement(comment);
        commentList.appendChild(commentElement);
    });

    postElement.appendChild(interactionsDiv);
    postElement.appendChild(commentsSection);
    toggleCommentsButton.addEventListener('click', () => {
        if (commentsSection.style.display === "none") {
            commentsSection.style.display = "block";
            toggleCommentsText.textContent = "Hide Comments";
        } else {
            commentsSection.style.display = "none";
            toggleCommentsText.textContent = "Show Comments";
        }});
    CommentHandler(commentButton, commentTextArea, commentList, post.ID);
    return postElement;
}