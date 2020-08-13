<template>
  <div
    :class="containerClasses"
    v-click-outside="clickOutside"
  >
    <v-text-field
      flat
      hide-details
      @focus="focus"
      v-model="internalInputValue"
      ref="input"
      solo-inverted
      :loading="loading"
      @keydown.enter.prevent.stop="onEnter"
      @keydown.tab.prevent.stop="onTab"
      @keydown.backspace="onBackspace"
      @keydown.up.prevent.stop="onUp"
      @keydown.down.prevent.stop="onDown"
      @keydown.esc.prevent.stop="onEsc"
      dense
      :class="$style.input"
    >
      <template #prepend-inner>
        <div class="d-flex align-center">
          <v-icon :color="isFocused ? 'primary' : ''">{{ mdiMagnify }}</v-icon>
          <template v-for="[type, value] in searchParams">
            <span :key="`${type}:${value}`" :class="$style.search">
              <span :class="$style.searchType">{{ type }}</span>
              <span :class="$style.searchValue">{{ value }} </span>
            </span>
          </template>
          <span v-if="hasValidSearchType" :class="$style.currentSearchType">{{ searchType }}:</span>
        </div>
      </template>
      <template #label>
        <span :class="$style.inputLabel">Search</span>
      </template>
      <div>NwefawefO input</div>
    </v-text-field>

    <v-menu
      :value="isMenuActive && menuCanShow"
      :close-on-click="false"
      :close-on-content-click="false"
      :open-on-click="false"
      :activator="$refs.input"
      :attach="$el"
      max-height="380"
      ref="menu"
      disable-keys
      offset-y
      offset-overflow
      :nudge-bottom="2"
    >
      <v-list class="v-select-list" tabindex="-1" dense>
        <v-list-item
          v-for="(item, n) in listData"
          :key="getItemID(item)"
          @click="onItemClick(item, n)"
          :class="{'v-list-item--highlighted': n == selectIndex}"
        >
          <v-list-item-content>
            <v-list-item-title v-if="isSearchingType">{{ item.text }}</v-list-item-title>
            <v-list-item-title v-if="isSearchingTag">
              <tag-chip :tag="item"></tag-chip>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import { mdiMagnify } from '@mdi/js'
import debounce from 'lodash.debounce'
import TagChip from '@/components/tag/TagChip'
import { SearchParams } from '@/utils/search'

export default {
  name: 'app-search',
  components: {
    TagChip,
  },
  data () {
    return {
      loading: false,
      isFocused: false,
      isMenuActive: false,
      inputValue: '',
      tags: [],
      selectIndex: -1,
      searchParams: new SearchParams(),

      mdiMagnify,
    }
  },
  created () {
    this.restoreSearchParams()
    document.addEventListener('keypress', this.onDocKeypress)
  },
  beforeDestroy () {
    document.removeEventListener('keypress', this.onDocKeypress)
  },
  computed: {
    ...mapState(['client']),
    hasSearchType () {
      return this.searchSeparatorIndex > -1
    },
    internalInputValue: {
      get () {
        if (this.hasValidSearchType) {
          return this.searchInput
        }
        return this.inputValue
      },
      set (val) {
        if (this.hasValidSearchType) {
          this.inputValue = `${this.searchType}:${val}`
        } else {
          this.inputValue = val
        }
      },
    },
    storedSearch: {
      get () {
        return this.$store.state.search
      },
      set (val) {
        this.$store.commit('SET_SEARCH', val)
      },
    },
    shouldSuggestSearchType () {
      return !this.hasSearchType && this.inputValue.length > 0
    },
    isSearchingType () {
      return this.shouldSuggestSearchType
    },
    isValidSearchType () {
      let inputValue = this.hasSearchType
        ? this.searchType
        : this.inputValue
      return this.allowedTypes.findIndex(({ name }) => name == inputValue) > -1
    },
    hasValidSearchType () {
      return this.hasSearchType && this.isValidSearchType
    },
    isSearchingTag () {
      return this.searchType == 'tag'
    },
    menuCanShow () {
      return !this.isEmptyList && (this.isSearchingType || this.isSearchingTag)
    },
    searchSeparatorIndex () {
      return this.inputValue.indexOf(':')
    },
    searchType () {
      if (!this.hasSearchType) {
        return null
      }
      return this.inputValue.slice(0, this.searchSeparatorIndex)
    },
    searchInput () {
      if (!this.hasSearchType) {
        return null
      }
      return this.inputValue.slice(this.searchSeparatorIndex + 1)
    },
    allowedTypes () {
      return [
        { name: 'tag', text: 'Tag', alias: ['tag', '#'] },
      ]
    },
    isEmptyList () {
      return this.listData.length == 0
    },
    filteredTags () {
      const selectedTags = this.searchParams.getValues('tag')
      return this.tags.filter(tag => !selectedTags.includes(`/${tag}`))
    },
    listData () {
      if (this.shouldSuggestSearchType) {
        return this.allowedTypes.filter(type => {
          return type.alias.some(alias => {
            return alias.includes(this.inputValue)
          })
        })
      }
      if (this.isSearchingTag) {
        return this.filteredTags
      }
      return []
    },
    isSelected () {
      return this.selectIndex > -1
    },
    clickOutside () {
      return {
        handler: this.blur,
      }
    },
    containerClasses () {
      return {
        [this.$style.container]: true,
        [this.$style.containerExpanded]: this.isFocused,
      }
    },
  },
  methods: {
    focus () {
      this.isFocused = true
      this.isMenuActive = true
      this.$refs.input.focus()
    },
    blur () {
      this.isFocused = false
      this.isMenuActive = false
      this.$refs.input.blur()
    },
    getItemID (item) {
      if (this.isSearchingType) {
        return item.name
      }
      if (this.isSearchingTag) {
        return item.id
      }
    },
    search () {
      this.storedSearch = this.searchParams.clone()
      if (this.$route.query.q != this.storedSearch.encode()) {
        this.$router.push({
          path: `/pin`, 
          query: { q: this.storedSearch.encode() },
        })
      }
    },
    restoreSearchParams () {
      if (this.$route.path != '/pin') {
        return
      }

      const params = new SearchParams()
      params.parse(this.$route.query.q)

      this.storedSearch = params.clone()
      this.searchParams = params
      this.inputValue = params.getText()
    },
    debounceFetchData: debounce(async function () {
      if (this.isSearchingTag) {
        await this.fetchTags()
      }
    }, 200),
    async fetchTags () {
      try {
        this.loading = true
        const res = await this.client.listTags({ q: this.searchInput })
        this.tags = res.data.map(tag => tag.name)
      } catch (e) {
        //
      } finally {
        this.loading = false
      }
    },
    setType (index) {
      const { name } = this.listData[index]
      this.inputValue = name
    },
    setTag (index) {
      this.internalInputValue = '/' + this.listData[index]
    },
    selectNextItem () {
      const item = this.listData[this.selectIndex + 1]

      if (!item) {
        if (!this.listData.length) {
          return
        }

        this.selectIndex = -1
        this.selectNextItem()
        return
      }

      this.selectIndex++
    },
    selectPrevItem () {
      const item = this.listData[this.selectIndex - 1]

      if (!item) {
        if (!this.listData.length) {
          return
        }

        this.selectIndex = this.listData.length
        this.selectPrevItem()
        return
      }

      this.selectIndex--
    },
    getSelectIndex () {
      let index = this.selectIndex
      index = Math.max(0, index)
      return index
    },
    onDocKeypress (e) {
      if (e.key == '/') {
        if (this.$el.contains(e.target)) {
          return
        }
        if (e.target.tagName == 'INPUT') {
          return
        }
        e.preventDefault()
        e.stopPropagation()
        this.focus()
      }
    },
    onEnter () {
      // When searching types with select.
      if (this.isSearchingType &&
        this.isValidSearchType &&
        this.isSelected
      ) {
        this.setType(this.getSelectIndex())
        this.inputValue += ':'
      // When searching tags.
      } else if (this.hasValidSearchType && this.isSearchingTag) {
        if (this.isSelected) {
          this.setTag(this.getSelectIndex())
        }
        if (this.searchInput == '') {
          return
        }

        this.searchParams.add(this.searchType, this.searchInput)
        this.inputValue = ''
      // Submit search to pin view.
      } else {
        this.searchParams.setText(this.inputValue)
        this.search()
      }
    },
    onTab () {
      // Skip if data is empty.
      if (this.isEmptyList) {
        // no-op
      // When searching types.
      } else if (this.isSearchingType) {
        this.setType(this.getSelectIndex())
        this.$nextTick(() => {
          this.selectIndex = 0
        })
      // When searching tags.
      } else if (this.isSearchingTag) {
        this.setTag(this.getSelectIndex())
      }
    },
    onBackspace (e) {
      // When empty input on searching type's items.
      if (this.hasSearchType && !this.searchInput) {
        e.preventDefault()
        this.inputValue = this.inputValue.slice(0, -1)
      // When input is empty and has filter.
      } else if (!this.inputValue && !this.searchParams.isEmptyParam()) {
        e.preventDefault()
        const search = this.searchParams.pop()
        this.inputValue = search.join(':')
      }
    },
    onItemClick (item, n) {
      if (this.isSearchingType) {
        this.setType(n)
        this.inputValue += ':'
        this.$refs.input.focus()
      }
    },
    onUp () {
      this.selectPrevItem()
    },
    onDown () {
      this.selectNextItem()
    },
    onEsc () {
      this.blur()
    },
  },
  watch: {
    inputValue () {
      this.selectIndex = -1
    },
    searchInput () {
      if (!this.loading) {
        this.debounceFetchData()
      }
    },
    '$route' () {
      this.restoreSearchParams()
    },
  },
}
</script>

<style lang="scss" module>
.container {
  position: relative;
  width: 300px;
  transition: all .2s ease;
}

.containerExpanded {
  width: 800px;
}

.input input,
.inputLabel {
  font-size: 14px;
}

$search-border-radius: 2px;
$search-height: 20px;
$search-gutter: 4px;

.search {
  font-size: 12px;
  display: inline-flex;
  align-items: center;
  background-color: #e0e0e0;
  height: $search-height;
  padding: 0 $search-gutter;
  border-radius: $search-border-radius;
  margin: 2px;
}

.searchType {
  display: block;
  background-color: black;
  line-height: $search-height;
  padding: 0 $search-gutter;
  margin-left: -1 * $search-gutter;
  margin-right: $search-gutter;
  border-top-left-radius: $search-border-radius;
  border-bottom-left-radius: $search-border-radius;
  color: #f0f0f0;
}

.searchValue {
  color: #000000de;
  white-space: nowrap;
}

.currentSearchType {
  font-size: 14px;
  margin-left: 4px;
}
</style>
