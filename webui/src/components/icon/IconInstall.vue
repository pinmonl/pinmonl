<template>
  <svg v-show="false">
    <symbol v-for="(d, k) in symbols" :key="k" viewBox="0 0 24 24" :id="k">
      <path :d="d" />
    </symbol>
  </svg>
</template>

<script>
import {
  mdiAccount,
  mdiAccountCircle,
  mdiBookmark,
  mdiClose,
  mdiContentSave,
  mdiCubeOutline,
  mdiDocker,
  mdiEye,
  mdiGithubCircle,
  mdiGitlab,
  mdiKubernetes,
  mdiLink,
  mdiMagnify,
  mdiPencil,
  mdiPlus,
  mdiPound,
  mdiSourceFork,
  mdiStar,
  mdiTag,
  mdiTagMultiple,
  mdiUndoVariant,
} from '@mdi/js'

export default {
  computed: {
    symbols () {
      return { ...this.originals, ...this.alias }
    },
    alias () {
      return {
        add: mdiPlus,
        cancel: mdiUndoVariant,
        edit: mdiPencil,
        github: mdiGithubCircle,
        save: mdiContentSave,
      }
    },
    originals () {
      return {
        account: mdiAccount,
        accountCircle: mdiAccountCircle,
        bookmark: mdiBookmark,
        close: mdiClose,
        contentSave: mdiContentSave,
        cubeOutline: mdiCubeOutline,
        docker: mdiDocker,
        eye: mdiEye,
        githubCircle: mdiGithubCircle,
        gitlab: mdiGitlab,
        kubernetes: mdiKubernetes,
        link: mdiLink,
        magnify: mdiMagnify,
        pencil: mdiPencil,
        plus: mdiPlus,
        pound: mdiPound,
        sourceFork: mdiSourceFork,
        star: mdiStar,
        tag: mdiTag,
        tagMultiple: mdiTagMultiple,
        undoVariant: mdiUndoVariant,
      }
    },
  },
}
</script>
