/** 多语言翻译定义 */

export type Locale = 'zh' | 'en'

/** 翻译文案结构 */
export interface Messages {
  // 页面标题区域
  pageTitle: string
  pageSubtitle: string

  // 上传区域
  uploadDropText: string
  uploadPasteHint: string
  uploadFormatHint: string
  uploadReleaseHint: string
  uploading: string

  // 上传结果
  uploadResults: string
  clearAll: string
  deleteImage: string
  copiedToClipboard: string
  imageDeleted: string
  deleteImageFailed: string

  // 上传提示消息
  uploadSuccess: (count: number) => string
  uploadFailed: string
  fileTooLarge: (name: string, maxMB: number) => string
  formatNotAllowed: (name: string) => string
  tooManyFiles: (max: number) => string

  // 页脚
  footerText: string

  // 主题
  themeLabel: string

  // 语言
  langLabel: string

  // 水印提示
  watermarkHint: string

  // 过期时间
  expiryLabel: string
  expiryNever: string
  expiryHour: (h: number) => string
  expiryDay: (d: number) => string
  expired: string
  expiresIn: string
}

/** 中文翻译 */
const zh: Messages = {
  pageTitle: '上传 & 分享图片',
  pageSubtitle: '快速、免费、无需登录，即刻获取分享链接。',

  uploadDropText: '拖拽图片到此处，或点击选择文件',
  uploadPasteHint: '也可使用 Ctrl+V 粘贴剪贴板图片',
  uploadFormatHint: '支持 JPG、PNG、GIF、WebP、SVG，单张最大 5MB',
  uploadReleaseHint: '松开鼠标即可上传',
  uploading: '上传中...',

  uploadResults: '上传结果',
  clearAll: '清空全部',
  deleteImage: '删除图片',
  copiedToClipboard: '已复制到剪贴板',
  imageDeleted: '图片已删除',
  deleteImageFailed: '删除图片失败',

  uploadSuccess: (count) => `成功上传 ${count} 张图片`,
  uploadFailed: '上传失败',
  fileTooLarge: (name, maxMB) => `${name} 超过最大限制 (${maxMB}MB)`,
  formatNotAllowed: (name) => `${name} 格式不被支持`,
  tooManyFiles: (max) => `单次最多上传 ${max} 张`,

  footerText: 'doyo-img - 轻量级图床服务',

  themeLabel: '主题',
  langLabel: '语言',

  watermarkHint: '上传的图片将添加水印',

  expiryLabel: '有效期',
  expiryNever: '永不过期',
  expiryHour: (h) => `${h} 小时`,
  expiryDay: (d) => `${d} 天`,
  expired: '已过期',
  expiresIn: '剩余有效期',
}

/** 英文翻译 */
const en: Messages = {
  pageTitle: 'Upload & Share Images',
  pageSubtitle: 'Fast, free, and no login required. Get shareable links instantly.',

  uploadDropText: 'Drop images here, or click to select',
  uploadPasteHint: 'Paste from clipboard with Ctrl+V',
  uploadFormatHint: 'Supports JPG, PNG, GIF, WebP, SVG \u00b7 Max 5MB per file',
  uploadReleaseHint: 'Release to upload',
  uploading: 'Uploading...',

  uploadResults: 'Upload Results',
  clearAll: 'Clear all',
  deleteImage: 'Delete image',
  copiedToClipboard: 'Copied to clipboard',
  imageDeleted: 'Image deleted',
  deleteImageFailed: 'Failed to delete image',

  uploadSuccess: (count) => `${count} image(s) uploaded successfully`,
  uploadFailed: 'Upload failed',
  fileTooLarge: (name, maxMB) => `${name} exceeds max size (${maxMB}MB)`,
  formatNotAllowed: (name) => `${name} is not an allowed format`,
  tooManyFiles: (max) => `Maximum ${max} files per upload`,

  footerText: 'doyo-img - Lightweight Image Hosting',

  themeLabel: 'Theme',
  langLabel: 'Language',

  watermarkHint: 'Uploaded images will be watermarked',

  expiryLabel: 'Expiry',
  expiryNever: 'Never',
  expiryHour: (h) => `${h}h`,
  expiryDay: (d) => `${d}d`,
  expired: 'Expired',
  expiresIn: 'Expires in',
}

/** 所有语言映射 */
const messages: Record<Locale, Messages> = { zh, en }

/** 根据语言代码获取翻译文案 */
export function getMessages(locale: Locale): Messages {
  return messages[locale]
}
