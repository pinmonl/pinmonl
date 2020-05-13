<template>
  <Anchor :class="$style.bmcButton" external to="https://www.buymeacoffee.com/pinmonl">
    <BMCLogo />
    <span :class="$style.bmcText">Buy me a coffee</span>
  </Anchor>
</template>

<script>
import BMCLogo from '@/assets/bmc-new-btn-logo.svg'

export default {
  components: {
    BMCLogo,
  },
}
</script>

<style lang="scss" module>
.bmcButton svg {
  height: 35px;
  width: 34px;
  margin-bottom: 1px;
  box-shadow: none;
  border: none;
  vertical-align: middle;
}
.bmcButton {
  padding: 7px 10px;
  height: 51px;
  line-height: 35px;
  min-width: 217px;
  text-decoration: none;
  display: inline-flex;
  color: #ffffff;
  background-color: #ff813f;
  border-radius: 5px;
  border: 1px solid transparent;
  font-size: 20px;
  letter-spacing: 0.6px;
  // box-shadow: 0px 1px 2px rgba(190, 190, 190, 0.5);
  // -webkit-box-shadow: 0px 1px 2px 2px rgba(190, 190, 190, 0.5);
  margin: 0 auto;
  font-family: "Lato", sans-serif;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
}
.bmcText {
  margin-left: 15px;
  font-size: 19px;
}
</style>
