<template>
  <div>
    <div>
      <md-empty-state
        md-rounded
        md-icon="highlight_off"
        md-label="Tivemos um erro ao coletar os dados. Veja o erro abaixo."
      >
      </md-empty-state>
    </div>

    <div v-show="this.executorLog.cmd != ''" class="errorLog">
      <b>Erro no comando: </b>
      <br />
      <textarea wrap="soft" class="textArea">
 {{ this.executorLog.cmd }} </textarea
      >
      <br />
      <b> Saída de erro: </b>
      <br />
      <textarea class="textArea">  {{ this.executorLog.err }}</textarea>
      <br />
      <b v-show="this.executorLog.stdout != ''"> Saída padrão: </b>

      <textarea class="textArea" v-show="this.executorLog.stdout != ''">
        {{ this.executorLog.stdout }}
      </textarea>
      <b> Variáveis de ambiente (env): </b>
      <br />
      <textarea rows="20" class="textArea">
  {{ this.executorLog.env }}</textarea
      >
    </div>
    <div
      style="text-align: center; margin-bottom: 10px"
      v-show="this.executorLog.cmd != ''"
    >
      <h5>
        Acha que tem algo errado? Por favor entre em contato conosco abrindo uma
        issue.
      </h5>
      <a
        style="padding-bottom: 15px; font-size: 16px;  font-weight: normal;"
        href="https://github.com/dadosjusbr/coletores/issues/new"
        target="_blank"
      >
        Abra uma issue aqui</a
      >
    </div>
  </div>
</template>

<script>
export default {
  name: "errorCollectingDataPage",
  props: {
    executorLog: {
      type: Object,
      default: {},
    },
  },
};
</script>

<style >
.errorLog {
  padding-left: 5px;
  padding-top: 5px;
  padding-right: 5px;
  padding-bottom: 5px;
  text-align: center;
  margin-bottom: 5px;
}
.textArea {
  border: 1px solid #2ab38b;
  width: 80%;
}
h5 {
  margin-bottom: 0px;
}
a {
  color: black;
  font-size: 1.4em;
  font-weight: bold;
}
</style>
