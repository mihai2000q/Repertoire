import { mantineRender } from '../../../test-utils.tsx'
import ImageDropzoneWithPreview from './ImageDropzoneWithPreview.tsx'
import {screen} from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Image Dropzone With Preview', () => {
  it('should show dropzone when image is null', () => {
    mantineRender(<ImageDropzoneWithPreview image={null} setImage={() => {}} />)

    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
  })

  it('should show preview when an image is uploaded on the dropzone', async () => {
    const user = userEvent.setup()

    const image = new File([''], 'image.png', { type: 'image/png' })
    const setImage = vitest.fn()

    const { rerender } = mantineRender(
      <ImageDropzoneWithPreview image={null} setImage={setImage} />
    )

    await user.upload(screen.getByTestId('image-dropzone-input'), image)

    expect(setImage).toHaveBeenCalledOnce()
    expect(setImage).toHaveBeenCalledWith(image)

    rerender(<ImageDropzoneWithPreview image={image} setImage={setImage} />)

    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'image-options' })).toBeInTheDocument()
  })

  it('should display menu on click on preview', async () => {
    const user = userEvent.setup()

    const image = new File([''], 'image.png', { type: 'image/png' })

    mantineRender(<ImageDropzoneWithPreview image={image} setImage={() => {}} />)

    await user.click(screen.getByRole('button', { name: 'image-options' }))
    expect(screen.getByRole('menuitem', { name: /upload image/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove image/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should re-upload image', async () => {
      const user = userEvent.setup()

      const image = new File([''], 'image.png', { type: 'image/png' })
      const newImage = new File([''], 'new-image.png', { type: 'image/png' })
      const setImage = vitest.fn()

      mantineRender(<ImageDropzoneWithPreview image={image} setImage={setImage} />)

      await user.click(screen.getByRole('button', { name: 'image-options' }))
      await user.upload(screen.getByTestId('upload-image-input'), newImage)

      expect(setImage).toHaveBeenCalledOnce()
      expect(setImage).toHaveBeenCalledWith(newImage)
    })

    it('should remove image and show the dropzone', async () => {
      const user = userEvent.setup()

      const image = new File([''], 'image.png', { type: 'image/png' })
      const setImage = vitest.fn()

      const { rerender } = mantineRender(<ImageDropzoneWithPreview image={image} setImage={setImage} />)

      await user.click(screen.getByRole('button', { name: 'image-options' }))
      await user.click(screen.getByRole('menuitem', { name: /remove image/i }))

      expect(setImage).toHaveBeenCalledOnce()
      expect(setImage).toHaveBeenCalledWith(null)

      rerender(<ImageDropzoneWithPreview image={null} setImage={setImage} />)

      expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    })
  })
})
