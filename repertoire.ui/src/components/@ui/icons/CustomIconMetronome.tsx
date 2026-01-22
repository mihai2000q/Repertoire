import { SVGProps } from 'react'

interface CustomIconMetronomeProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

function CustomIconMetronome({ color, size = 24, strokeWidth = 2, ...props }: CustomIconMetronomeProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink={'http://www.w3.org/1999/xlink'}
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      strokeWidth={strokeWidth}
      strokeLinecap="round"
      strokeLinejoin="round"
      viewBox="0 0 512 512"
      xmlSpace={'preserve'}
      {...props}
    >
      <g>
        <g>
          <path
            d="M451.555,144.609c-7.08-7.08-18.56-7.08-25.64,0l-32.765,32.766l-1.238-1.238c-7.08-7.08-18.56-7.08-25.64,0
			c-6.174,6.172-8.393,17.247,1.238,26.879c-6.777,6.777-148.384,148.384-154.836,154.836c18.276,0,33.653,0,51.281,0
			c17.721-17.72,127.206-127.206,129.196-129.196l1.238,1.238c7.081,7.081,18.56,7.08,25.64,0c7.081-7.08,7.081-18.56,0-25.64
			l-1.238-1.238l32.766-32.765C458.636,163.169,458.636,151.689,451.555,144.609z"
          />
        </g>
      </g>
      <g>
        <g>
          <polygon points="362.468,307.361 311.977,357.851 369.122,357.851 		" />
        </g>
      </g>
      <g>
        <g>
          <path
            d="M342.038,152.357l-6.018-45.663c-0.628-4.765-2.892-9.166-6.406-12.448L234.461,5.375c-7.673-7.166-19.585-7.166-27.258,0
			l-95.152,88.872c-3.513,3.281-5.777,7.681-6.406,12.448L72.542,357.851h92.668l38.642-39.204V124.969
			c0-9.377,7.602-16.979,16.979-16.979c9.377,0,16.979,7.602,16.979,16.979v159.72l89.552-89.552
			C325.519,179.456,330.693,163.836,342.038,152.357z"
          />
        </g>
      </g>
      <g>
        <g>
          <path
            d="M386.348,488.549l-12.75-96.74H68.067l-12.751,96.74C53.679,500.971,63.359,512,75.878,512h289.907
			C378.315,512,387.985,500.962,386.348,488.549z"
          />
        </g>
      </g>
    </svg>
  )
}

export default CustomIconMetronome
