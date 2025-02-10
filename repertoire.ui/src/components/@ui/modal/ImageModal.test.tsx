import { mantineRender } from '../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import ImageModal from './ImageModal.tsx'

describe('Image Modal', () => {
  it('should render', () => {
    const title = 'Image Modal'
    const image = 'something.png'

    mantineRender(<ImageModal opened={true} onClose={() => {}} title={title} image={image} />)

    expect(screen.getByRole('dialog', { name: title + '-image' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: title })).toBeInTheDocument()
    expect(screen.getByRole('button')).toBeInTheDocument() // close button
  })
})
