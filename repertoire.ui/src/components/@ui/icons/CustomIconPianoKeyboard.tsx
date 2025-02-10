import { forwardRef, SVGProps } from 'react'

interface CustomIconPianoKeyboardProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconPianoKeyboard = forwardRef<SVGSVGElement, CustomIconPianoKeyboardProps>(
  ({ color, size = 24, ...props }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink={'http://www.w3.org/1999/xlink'}
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      version="1.1"
      viewBox="0 0 511.999 511.999"
      xmlSpace={'preserve'}
      {...props}
    >
      <g>
        <g>
          <path
            d="M478.168,91.547H33.831C15.176,91.547,0,106.723,0,125.379v261.243c0,18.655,15.176,33.831,33.831,33.831h444.337
			c18.655,0,33.831-15.176,33.831-33.831V125.379C512,106.723,496.823,91.547,478.168,91.547z M79.779,386.993H33.831
			c-0.205,0-0.371-0.166-0.371-0.371V175.696h35.166v100.992c0,4.466,3.62,8.086,8.086,8.086h3.067V386.993z M161.197,386.993
			h-47.959V284.775h2.509c4.466,0,8.086-3.62,8.086-8.086V175.696h26.489v100.992c0,4.466,3.62,8.086,8.086,8.086h2.788V386.993z
			 M239.27,386.993h-44.613V284.775h2.788c4.466,0,8.086-3.62,8.086-8.086V175.696h22.864v100.992c0,4.466,3.62,8.086,8.086,8.086
			h2.788V386.993z M320.688,386.993h-47.959V284.775h2.788c4.466,0,8.086-3.62,8.086-8.086V175.696h26.21v100.992
			c0,4.466,3.62,8.086,8.086,8.086h2.788V386.993z M398.76,386.993h-44.613V284.775h2.788c4.466,0,8.086-3.62,8.086-8.086V175.696
			h22.864v100.992c0,4.466,3.62,8.086,8.086,8.086h2.788V386.993z M478.539,386.621c0.001,0.205-0.166,0.371-0.371,0.371H432.22
			V284.775h2.788c4.466,0,8.086-3.62,8.086-8.086V175.696h35.445V386.621z"
          />
        </g>
      </g>
    </svg>
  )
)

CustomIconPianoKeyboard.displayName = 'CustomIconPianoKeyboard'

export default CustomIconPianoKeyboard
