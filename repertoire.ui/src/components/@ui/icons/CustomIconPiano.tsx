import { forwardRef, SVGProps } from 'react'

interface CustomIconPianoProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconPiano = forwardRef<SVGSVGElement, CustomIconPianoProps>(
  ({ color, size = 24, ...props }, ref) => (
    <svg
      ref={ref}
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      version="1.1"
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink={'http://www.w3.org/1999/xlink'}
      viewBox="0 0 512 512"
      xmlSpace="preserve"
      {...props}
    >
      <g>
        <g>
          <g>
            <path
              d="M432.149,225.323c-58.005-7.296-81.664-31.232-81.664-82.624v-4.032C350.485,62.208,286.379,0,207.616,0h-54.101
				C74.752,0,10.667,62.208,10.667,138.667V320h490.667v-18.496C501.333,263.019,471.595,230.251,432.149,225.323z"
            />
            <path
              d="M437.333,426.667c0,11.776-9.557,21.333-21.333,21.333s-21.333-9.557-21.333-21.333v-64h-21.333v64
				c0,11.776-9.557,21.333-21.333,21.333s-21.333-9.557-21.333-21.333v-64h-21.333v64c0,11.776-9.557,21.333-21.333,21.333
				s-21.333-9.557-21.333-21.333v-64h-21.333v64c0,11.776-9.557,21.333-21.333,21.333s-21.333-9.557-21.333-21.333v-64h-21.333v64
				c0,11.776-9.557,21.333-21.333,21.333s-21.333-9.557-21.333-21.333v-64h-21.333v64C117.333,438.443,107.776,448,96,448
				s-21.333-9.557-21.333-21.333v-64h-64V512h490.667V362.667h-64V426.667z"
            />
          </g>
        </g>
      </g>
    </svg>
  )
)

CustomIconPiano.displayName = 'CustomIconPiano'

export default CustomIconPiano
