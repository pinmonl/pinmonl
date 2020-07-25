import {
  mdiGithub,
} from '@mdi/js'

export function getProvider (provider) {
  const map = {
    'github': { label: 'GitHub', icon: mdiGithub },
  }
  return map[provider]
}