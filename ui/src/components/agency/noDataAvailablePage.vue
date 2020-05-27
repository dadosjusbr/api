<template>
  <div>
    <md-empty-state
      md-rounded
      md-icon="highlight_off"
      md-label="Talvez o órgão ainda não tenha disponibilizado os dados ou o dadosjusbr não tentou realizar a coleta."
      md-description="Acha que tem algo errado? Por favor entre em contato conosco abrindo uma issue."
    >
      <a
        style="padding-bottom: 15px; font-size: 16px;  font-weight: normal;"
        href="https://github.com/dadosjusbr/coletores/issues/new"
        target="_blank"
      >
        Abra uma issue aqui</a
      >
    </md-empty-state>
  </div>
</template>

<script>
export default {
  name: "noDataAvailablePage",
};
</script>

<style></style>
