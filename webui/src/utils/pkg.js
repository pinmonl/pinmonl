import {
  mdiDocker,
  mdiGithub,
  mdiNpm,
  mdiYoutube,
} from '@mdi/js'

export function getProvider (provider) {
  const map = {
    'docker': { label: 'Docker', icon: mdiDocker },
    'github': { label: 'GitHub', icon: mdiGithub },
    'npm': { label: 'NPM', icon: mdiNpm },
    'youtube': { label: 'YouTube', icon: mdiYoutube },
  }
  return map[provider]
}
