import type { LinkFormat } from '../types'

export function generateLink(url: string, format: LinkFormat, name?: string): string {
  const alt = name || 'image'
  switch (format) {
    case 'url':
      return url
    case 'markdown':
      return `![${alt}](${url})`
    case 'html':
      return `<img src="${url}" alt="${alt}" />`
    case 'bbcode':
      return `[img]${url}[/img]`
    default:
      return url
  }
}

export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + units[i]
}

export function isAllowedFormat(file: File, allowedFormats: string[]): boolean {
  const ext = file.name.split('.').pop()?.toLowerCase() || ''
  const mimeMap: Record<string, string[]> = {
    jpg: ['image/jpeg'],
    jpeg: ['image/jpeg'],
    png: ['image/png'],
    gif: ['image/gif'],
    webp: ['image/webp'],
    svg: ['image/svg+xml'],
  }
  return allowedFormats.some((fmt) => {
    if (ext === fmt) return true
    const mimes = mimeMap[fmt]
    return mimes?.includes(file.type)
  })
}
