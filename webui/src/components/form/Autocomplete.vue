<template>
  <div :class="$style.wrapper">
    <div :class="$style.container">
      <slot name="input" v-bind="slotBinds()"></slot>
    </div>
    <div :class="[$style.options, optionsClass]" v-if="show" ref="options">
      <div
        v-for="(opt, n) in options"
        :key="getValue(opt)"
        ref="option"
        :class="[$style.optionWrapper, optionWrapperClass]"
      >
        <slot
          name="option"
          v-bind="slotBinds({
            option: opt,
            n: n,
            move: (v) => move(n, v),
            active: isActive(n),
            selected: () => isSelected(n),
            class: [$style.option, computeOptionClass(opt, n)],
          })"
        >
          <div
            :class="[$style.option, computeOptionClass(opt, n), optionClass]"
            @mouseover="move(n)"
            @click="toggle(n)"
          >
            {{ getLabel(opt) }}
          </div>
        </slot>
      </div>
      <slot
        name="create"
        v-if="creatable"
        v-bind="slotBinds({
          class: [$style.option, computeOptionClass(null, len)],
          move: () => move(len),
        })"
      />
      <slot name="after" />
    </div>
  </div>
</template>

<script>
export default {
  props: {
    options: {
      type: Array,
      default: () => ([]),
    },
    value: {
      type: [Array, Object],
      default: null,
    },
    valueBy: {
      type: [String, Function],
      required: true,
    },
    labelBy: {
      type: [String, Function],
      required: true,
    },
    show: {
      type: Boolean,
      default: false,
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    creatable: {
      type: Boolean,
      default: false,
    },
    optionsClass: [Array, Object, String],
    optionWrapperClass: [Array, Object, String],
    optionClass: [Array, Object, String],
  },
  data () {
    return {
      cursor: -1,
    }
  },
  computed: {
    len () {
      return this.options.length
    },
  },
  methods: {
    next () {
      this.cursor = Math.min(this.options.length, this.cursor + 1)
      this.scrollToCursor()
    },
    prev () {
      this.cursor = Math.max(-1, this.cursor - 1)
      this.scrollToCursor()
    },
    move (n, visible = false) {
      this.cursor = n
      if (visible) {
        this.scrollToCursor()
      }
    },
    select (n) {
      if (typeof this.options[n] == 'undefined') {
        return
      }
      const opt = this.options[n]
      const newValue = this.multiple
        ? [ ...(this.value || []), opt ]
        : opt
      this.$emit('input', newValue)
    },
    unselect (n) {
      if (typeof this.options[n] == 'undefined') {
        return
      }
      const newValue = this.multiple
        ? [ ...(this.value || []) ].filter(v => this.getValue(v) != this.getValue(this.options[n]))
        : null
      this.$emit('input', newValue)
    },
    toggle (n) {
      if (typeof this.options[n] != 'undefined') {
        return this.isSelected(this.options[n])
          ? this.unselect(n)
          : this.select(n)
      }
      if (this.creatable && this.len == this.cursor) {
        return this.$emit('create')
      }
    },
    autoSelect () {
      this.select(this.cursor)
    },
    autoUnselect () {
      this.unselect(this.cursor)
    },
    autoToggle () {
      this.toggle(this.cursor)
    },
    getValue (option) {
      return this.attrGetter(option, this.valueBy)
    },
    getLabel (option) {
      return this.attrGetter(option, this.labelBy)
    },
    attrGetter (option, getter) {
      if (typeof getter == 'string') {
        return option[getter]
      }
      if (typeof getter == 'function') {
        return getter(option)
      }
      return null
    },
    isSelected (option) {
      if (this.value == null) {
        return false
      }
      if (!this.multiple) {
        return this.getValue(this.value) == this.getValue(option)
      }
      return this.value
        .map(this.getValue)
        .includes(this.getValue(option))
    },
    isActive (option, n) {
      return this.cursor == n
    },
    computeOptionClass (option, n) {
      const cns = {
        [this.$style.option_active]: n == this.cursor,
      }
      if (option != null) {
        cns[this.$style.option_selected] = this.isSelected(option)
      }
      return cns
    },
    slotBinds (bindMore = {}) {
      return {
        next: this.next,
        prev: this.prev,
        moveTo: this.move,
        cursor: this.cursor,
        select: this.autoSelect,
        unselect: this.autoUnselect,
        toggle: this.autoToggle,
        selectAt: this.select,
        unselectAt: this.unselect,
        toggleAt: this.toggle,
        ...bindMore,
      }
    },
    scrollToCursor () {
      if (!this.show) {
        return
      }
      if (typeof this.options[this.cursor] == 'undefined') {
        return
      }
      const $con = this.$refs.options
      const $opt = this.$refs.option[this.cursor]
      const top = $con.scrollTop
      const bot = top + $con.clientHeight
      const oTop = $opt.offsetTop
      const oBot = oTop + $opt.clientHeight

      if (oTop < top) {
        $con.scrollTop = oTop
      }
      if (bot < oBot) {
        $con.scrollTop = oBot - $con.clientHeight
      }
    },
  },
  watch: {
    options (options) {
      if (options.length < this.cursor) {
        this.cursor = -1
      }
    },
  }
}
</script>

<style lang="scss" module>
.wrapper {
  @apply relative;
}

.container {
  @apply relative;
}

.options {
  @apply absolute;
  @apply w-full;
  top: 100%;
  @apply inset-x-0;
  margin-top: 1px;
  @apply shadow;
  @apply overflow-x-hidden;
  @apply overflow-y-auto;
  max-height: 300px;
  @apply bg-container;
  @apply z-50;
}

.option {
  @apply px-4;
  @apply py-1;
  @apply text-sm;

  &:hover {
    @apply cursor-pointer;
  }
}

.option_active {
  @apply bg-background;
}

.option_selected {
  @apply font-bold;
}
</style>
