import { forwardRef, SVGProps } from 'react'

interface CustomIconMusicNoteProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconMusicNote = forwardRef<SVGSVGElement, CustomIconMusicNoteProps>(
  ({ color, size = 24, ...props }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      strokeWidth={0}
      viewBox="-3 0 24 24"
      xmlSpace={'preserve'}
      {...props}
    >
      <path d="m18.07.169c-.148-.106-.333-.169-.532-.169-.111 0-.217.02-.316.055l.006-.002-11.077 3.938c-.361.131-.613.471-.613.869v.001 2.193.042 10.604c-.534-.295-1.169-.469-1.846-.471h-.001c-.043-.002-.093-.003-.143-.003-1.904 0-3.458 1.497-3.549 3.379v.008c.091 1.89 1.645 3.388 3.549 3.388.05 0 .1-.001.15-.003h-.007c.043.002.093.003.143.003 1.904 0 3.458-1.497 3.549-3.379v-.008-12.883l9.23-3.223v8.973c-.534-.294-1.17-.468-1.846-.47h-.001c-.043-.002-.094-.003-.144-.003-1.904 0-3.457 1.498-3.547 3.379v.008c.09 1.89 1.644 3.388 3.548 3.388.051 0 .101-.001.151-.003h-.007c.031.001.068.002.105.002 1.696 0 3.12-1.166 3.513-2.74l.005-.025c.042-.101.068-.217.069-.34v-15.754c0-.31-.153-.585-.388-.752l-.003-.002z" />
    </svg>
  )
)

CustomIconMusicNote.displayName = 'CustomIconMusicNote'

export default CustomIconMusicNote
