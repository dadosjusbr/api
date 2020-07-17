<template>
  <div class="footerContainer">
    <md-toolbar class="md-primary">
      <div class="logoContainer">
        <img src="../../assets/white_logo.png" />
      </div>
      <div class="contatcInfo">
        <h6>CONTATO</h6>
        <p style="margin-right: 15%">
          Alguma dica? Tem algum feedback geral? Se você tiver uma ideia que
          você gostaria de ver no DadosJusBr, envie-nos um e-mail
          (dadosjusbr@gmail.com) ou visite nosso
          <a style="color: white" href="https://github.com/dadosjusbr">github</a
          >.
        </p>
      </div>
    </md-toolbar>
  </div>
</template>

<script>
export default {
  name: "pageFooter",
};
</script>

<style lang="css" scoped>
.footerContainer {
  display: table;
  width: 100%;
  position: absolute;
  bottom: 0;
}
.md-toolbar {
  display: table-row;
  width: 100%;
  background-color: #3e5363;
  height: 15em;
}
.logoContainer {
  display: table-cell;
  width: 30%;
  text-align: center;
  vertical-align: middle;
}
img {
  max-width: 200px;
  max-height: 150px;
  width: auto;
  height: auto;
}
.contatcInfo {
  display: table-cell;
  text-align: left;
  color: white;
  font-size: 15px;
  vertical-align: middle;
}
@media only screen and (max-width: 780px) {
  .contatcInfo {
    padding-left: 10%;
  }
  .logoContainer{
    padding-left: 10%;
  }
}
@media only screen and (max-width: 600px) {
  .logoContainer {
    order: 2;
    position: absolute;
    bottom: 0;
    width: 100%;
    height: 40%;
    text-align: center;
    vertical-align: middle;
  }
  .contatcInfo {
    order: 1;
    padding-left: 10%;
    padding-top: 20%;
    text-align: left;
    color: white;
    font-size: 15px;
  }
  img {
    margin-top: 15%;
    vertical-align: middle;
    max-width: 200px;
    max-height: 150px;
    width: auto;
    height: auto;
  }
  .md-toolbar {
    flex-direction: column;
    display: flex;
    width: 100%;
    height: 23em;
  }
}
</style>
