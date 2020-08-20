import docker from 'simple-icons/icons/docker'
import git from 'simple-icons/icons/git'
import github from 'simple-icons/icons/github'
import npm from 'simple-icons/icons/npm'
import youtube from 'simple-icons/icons/youtube'

const icons = {
  docker,
  git,
  github,
  npm,
  youtube,
}

const getIcon = (name) => icons[name]

export {
  icons,
  getIcon,
}
