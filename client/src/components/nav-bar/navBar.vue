<template>
  <div class="topnav">
    <a class="active">DadosJus</a>
    <input type="text" placeholder="MPPB, Paraíba, RJ..." />
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
  background-color: #e9e9e9;
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
.topnav a:hover {
  background-color: #ddd;
  color: black;
}

/* Style the "active" element to highlight the current page */
.topnav a.active {
  background-color: #2b60de;
  color: white;
}

/* Style the search box inside the navigation bar */
.topnav input[type="text"] {
  float: right;
  padding: 6px;
  border: none;
  margin-top: 12px;
  margin-right: 16px;
  font-size: 17px;
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
}
</style>
