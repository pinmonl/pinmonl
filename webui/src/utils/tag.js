export const sanitizeTagName = name =>
  name.replace(/^\/+/, '').replace(/\/+$/, '')

export const absTagName = name => '/' + sanitizeTagName(name)

export const baseTagName = name => sanitizeTagName(name).replace(/^.*\//, '')
