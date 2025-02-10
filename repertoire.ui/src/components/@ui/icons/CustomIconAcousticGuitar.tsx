import { forwardRef, SVGProps } from 'react'

interface CustomIconAcousticGuitarProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconAcousticGuitar = forwardRef<SVGSVGElement, CustomIconAcousticGuitarProps>(
  ({ color, size = 24, ...props }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink={'http://www.w3.org/1999/xlink'}
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      version={'1.1'}
      viewBox="0 0 512 512"
      xmlSpace={'preserve'}
      {...props}
    >
      <g>
        <polygon points="274.685,80.761 237.3,80.761 233.499,236.663 278.5,236.663 	" />
        <rect x="234.178" width="43.636" height="64.48" />
        <rect x="213.09" y="1.725" width="12.869" height="20.998" />
        <rect x="213.09" y="41.75" width="12.869" height="20.998" />
        <rect x="286.034" y="1.725" width="12.868" height="20.998" />
        <rect x="286.034" y="41.75" width="12.868" height="20.998" />
        <path
          d="M256,266.987c-11.258,0-21.409,4.549-28.79,11.922c-7.373,7.38-11.921,17.532-11.921,28.79
		c0,11.258,4.548,21.409,11.921,28.782c7.381,7.373,17.532,11.922,28.79,11.922c11.259,0,21.41-4.549,28.79-11.922
		c7.374-7.374,11.922-17.525,11.922-28.782c0-11.258-4.549-21.41-11.922-28.79C277.409,271.536,267.258,266.987,256,266.987z
		 M256,338.482c-16.974,0-30.782-13.808-30.782-30.782c0-16.967,13.807-30.782,30.782-30.782c16.975,0,30.782,13.815,30.782,30.782
		C286.782,324.674,272.975,338.482,256,338.482z"
        />
        <path
          d="M383.769,358.93c-8.778-14.266-18.982-26.493-25.401-36.812c-2.992-4.808-4.412-10.258-4.419-16.44
		c-0.008-7.19,2.015-15.273,5.838-23.37c4.87-10.457,10.205-23.974,10.274-37.935c-0.008-32.682-26.485-59.167-59.168-59.167
		h-25.836l1.45,59.274h-61.015l1.443-59.274h-25.844c-32.682,0-59.16,26.485-59.168,59.167c0.069,13.968,5.42,27.493,10.282,37.935
		c3.823,8.106,5.854,16.196,5.846,23.37c-0.015,6.174-1.427,11.632-4.427,16.44c-6.411,10.319-16.616,22.547-25.394,36.804
		c-8.747,14.258-16.188,30.927-16.196,50.314c0,2.519,0.122,5.083,0.389,7.701c1.305,12.769,4.587,25.234,10.51,36.736
		c8.846,17.273,23.753,32.164,45.566,42.361C190.322,506.261,218.929,512,256,512c49.391,0,83.897-10.212,107.047-27.645
		c11.541-8.701,20.112-19.188,26.019-30.683c5.923-11.502,9.206-23.966,10.511-36.743c0.26-2.596,0.389-5.16,0.389-7.671
		C399.958,389.864,392.516,373.195,383.769,358.93z M289.216,421.997h-66.434v-10.289h66.434V421.997z M256,356.228
		c-26.806-0.008-48.527-21.73-48.527-48.528c0-26.806,21.722-48.528,48.527-48.536c26.806,0.008,48.528,21.73,48.528,48.536
		C304.528,334.498,282.805,356.22,256,356.228z"
        />
      </g>
    </svg>
  )
)

CustomIconAcousticGuitar.displayName = 'CustomIconAcousticGuitar'

export default CustomIconAcousticGuitar
