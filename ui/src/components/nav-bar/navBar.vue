<template>
  <div class="topnav">
    <router-link to="/">
      <a>
        <img src="../../assets/white_logo.png" class="active" />
      </a>
    </router-link>
  </div>
</template>

<script>
export default {
  name: "navBar"
};
</script>

<style>
/* Add a black background color to the top navigation bar */
.topnav {
  overflow: hidden;
  background-color: #182825;
}

/* Style the links inside the navigation bar */
.topnav a {
  font-family: "Open Sans Condensed", sans-serif;
  float: left;
  display: block;
  color: black;
  text-align: center;
  padding: 14px 16px;
  text-decoration: none;
  font-size: 25px;
}

/* Change the color of links on hover */
.topnav img:hover {
  background-color: black;
  color: black;
}

/* Style the "active" element to highlight the current page */
.topnav img.active {
  width: 15%;
  color: white;
}

/* When the screen is less than 600px wide, stack the links and the search field vertically instead of horizontally */
@media screen and (max-width: 600px) {
  .topnav a,
  .topnav input[type="text"] {
    float: none;
    display: block;
    text-align: left;
    width: 100%;
    margin: 0;
    padding: 14px;
  }
  .topnav input[type="text"] {
    border: 1px solid #ccc;
  }

  a:hover {
    background-color: yellow;
  }
}
</style>
