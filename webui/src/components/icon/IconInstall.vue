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
  mdiAlertCircleOutline,
  mdiBitbucket,
  mdiBookmark,
  mdiBookmarkPlus,
  mdiChevronRight,
  mdiClose,
  mdiContentSave,
  mdiCubeOutline,
  mdiDocker,
  mdiDotsVertical,
  mdiDownload,
  mdiEye,
  mdiFolderOutline,
  mdiGit,
  mdiGithub,
  mdiGitlab,
  mdiKubernetes,
  mdiLink,
  mdiMagnify,
  mdiNpm,
  mdiPencil,
  mdiPlus,
  mdiPound,
  mdiSourceFork,
  mdiStar,
  mdiTag,
  mdiTagMultiple,
  mdiUndoVariant,
} from '@mdi/js'

import helm from './helm.js'

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
        issue: mdiAlertCircleOutline,
        save: mdiContentSave,
      }
    },
    originals () {
      return {
        account: mdiAccount,
        accountCircle: mdiAccountCircle,
        alertCircleOutline: mdiAlertCircleOutline,
        bitbucket: mdiBitbucket,
        bookmark: mdiBookmark,
        bookmarkPlus: mdiBookmarkPlus,
        chevronRight: mdiChevronRight,
        close: mdiClose,
        contentSave: mdiContentSave,
        cubeOutline: mdiCubeOutline,
        docker: mdiDocker,
        dotsVertical: mdiDotsVertical,
        download: mdiDownload,
        eye: mdiEye,
        folderOutline: mdiFolderOutline,
        git: mdiGit,
        github: mdiGithub,
        gitlab: mdiGitlab,
        helm: helm,
        kubernetes: mdiKubernetes,
        link: mdiLink,
        magnify: mdiMagnify,
        npm: mdiNpm,
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
