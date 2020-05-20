export default () => {
  return {
    mounted () {
      document.addEventListener('keyup', this.handleKeyPress)
    },
    beforeDestroy () {
      document.removeEventListener('keyup', this.handleKeyPress)
    },
    computed: {
      shouldDisableKeys () {
        return this.$store.getters.globalSearch
      },
    },
    methods: {
      handleKeyPress () {
        // extends method
      },
    },
  }
}
