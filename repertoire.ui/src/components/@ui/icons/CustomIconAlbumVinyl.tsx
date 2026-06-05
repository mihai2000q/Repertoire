import { SVGProps } from 'react'

interface CustomIconAlbumVinylProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

function CustomIconAlbumVinyl({ color, size = 24, ...props }: CustomIconAlbumVinylProps) {
  return (
    <svg
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      version="1.1"
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink={'http://www.w3.org/1999/xlink'}
      viewBox="0 0 512 512"
      xmlSpace={'preserve'}
      {...props}
    >
      <g>
        <g>
          <path
            d="M256,0C114.842,0,0,114.84,0,256s114.842,256,256,256s256-114.84,256-256S397.158,0,256,0z M256,62.439v37.463
			c-36.3,0-69.747,12.455-96.29,33.317l-26.631-26.632C166.532,79.018,209.366,62.439,256,62.439z M99.902,256H62.439
			c0-46.634,16.578-89.469,44.149-122.921l26.63,26.63C112.359,186.253,99.902,219.699,99.902,256z M256,368.39
			c-62.071,0-112.39-50.318-112.39-112.39S193.929,143.61,256,143.61S368.39,193.928,368.39,256S318.071,368.39,256,368.39z"
          />
        </g>
      </g>
      <g>
        <g>
          <path
            d="M256,199.805c-30.986,0-56.195,25.209-56.195,56.195s25.209,56.195,56.195,56.195s56.195-25.209,56.195-56.195
			S286.986,199.805,256,199.805z M256,274.732c-10.329,0-18.732-8.403-18.732-18.732s8.403-18.732,18.732-18.732
			s18.732,8.403,18.732,18.732S266.329,274.732,256,274.732z"
          />
        </g>
      </g>
    </svg>
  )
}

export default CustomIconAlbumVinyl
