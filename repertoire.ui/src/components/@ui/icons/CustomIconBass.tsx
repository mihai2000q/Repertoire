import { forwardRef, SVGProps } from 'react'

interface CustomIconBassProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconBass = forwardRef<SVGSVGElement, CustomIconBassProps>(
  ({ color, size = 24, strokeWidth = 2, ...props }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink="http://www.w3.org/1999/xlink"
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      strokeWidth={2}
      version={'1.1'}
      viewBox="0 0 66.208 66.208"
      xmlSpace={'preserve'}
      {...props}
    >
      <path d="M65.51,1.579c-0.637-0.153-1.202-0.419-1.641-0.675c-0.596-0.348-1.345-0.235-1.85,0.236
	c-1.782,1.66-3.946,2.753-5.529,3.408c-0.719,0.298-1.161,1.058-0.996,1.818c0.263,1.212-0.708,2.259-1.054,2.587l-0.126,0.116l0,0
	l-30.11,27.876c-0.268,0.25-1.043,0.866-1.866,0.043c-1.038-1.038-1.224-4.378,3.486-9.089c0-0.062,1.266-0.809,0.457-1.619
	c-0.809-0.809-3.175,0.353-4.773,1.951c-4.71,4.711-4.677,8.106-7.097,10.459c-2.967,2.884-8.28,2.386-11.994,6.101
	c0,0-6.911,6.01,2.634,15.555c7.167,7.167,11.373,5.729,14.582,2.519c4.234-4.234,2.238-8.29,3.586-9.638
	c2.269-2.27,4.657-1.763,6.956-4.611c1.273-1.577,0.131-3.05-1.005-1.914c-0.939,0.939-3.422,2.135-4.75,0.807
	c-1.106-1.106-0.712-2.557,1.934-5.496l30.095-30.095c1.256-1.256,2.5-1.67,3.41-1.776c0.827-0.097,1.507-0.638,1.828-1.407
	c0.855-2.046,2.778-4.161,4.231-5.567C66.444,2.66,66.221,1.749,65.51,1.579z M15.874,56.741c-0.587,0.587-1.54,0.587-2.127,0
	l-5.64-5.64c-0.587-0.587-0.587-1.54,0-2.127c0.587-0.587,1.54-0.587,2.127,0l5.64,5.64C16.462,55.201,16.462,56.154,15.874,56.741z
	"/>    </svg>
  )
)

CustomIconBass.displayName = 'CustomIconBass'

export default CustomIconBass
