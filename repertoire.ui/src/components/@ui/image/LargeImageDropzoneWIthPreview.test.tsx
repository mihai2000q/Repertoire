import { mantineRender } from '../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import LargeImageDropzoneWithPreview from './LargeImageDropzoneWithPreview.tsx'

describe('Large Image Dropzone With Preview', () => {
  it('should show dropzone when image is null', () => {
    mantineRender(<LargeImageDropzoneWithPreview image={null} setImage={() => {}} />)

    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
  })

  it('should show preview when an image is uploaded on the dropzone', async () => {
    const user = userEvent.setup()

    const image = new File([''], 'image.png', { type: 'image/png' })
    const setImage = vitest.fn()

    const { rerender } = mantineRender(
      <LargeImageDropzoneWithPreview image={null} setImage={setImage} />
    )

    await user.upload(screen.getByTestId('image-dropzone-input'), image)

    expect(setImage).toHaveBeenCalledOnce()
    expect(setImage).toHaveBeenCalledWith(image)

    rerender(<LargeImageDropzoneWithPreview image={image} setImage={setImage} />)

    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()

    expect(screen.queryByRole('button', { name: 'reset-image' })).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'remove-image' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'upload-image' })).toBeInTheDocument()
  })

  it('should re-upload image', async () => {
    const user = userEvent.setup()

    const image = new File([''], 'image.png', { type: 'image/png' })
    const newImage = new File([''], 'new-image.png', { type: 'image/png' })
    const setImage = vitest.fn()

    mantineRender(<LargeImageDropzoneWithPreview image={image} setImage={setImage} />)

    await user.upload(screen.getByTestId('upload-image-input'), newImage)

    expect(setImage).toHaveBeenCalledOnce()
    expect(setImage).toHaveBeenCalledWith(newImage)
  })

  it('should remove image and show the dropzone', async () => {
    const user = userEvent.setup()

    const image = new File([''], 'image.png', { type: 'image/png' })
    const setImage = vitest.fn()

    const { rerender } = mantineRender(
      <LargeImageDropzoneWithPreview image={image} setImage={setImage} />
    )

    await user.click(screen.getByRole('button', { name: 'remove-image' }))

    expect(setImage).toHaveBeenCalledOnce()
    expect(setImage).toHaveBeenCalledWith(null)

    rerender(<LargeImageDropzoneWithPreview image={null} setImage={setImage} />)

    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
  })

  it('should reset image when default value is provided', async () => {
    const user = userEvent.setup()

    const defaultImage = 'default.png'
    const newImage = new File([''], 'image.png', { type: 'image/png' })
    const setImage = vitest.fn()

    const { rerender } = mantineRender(
      <LargeImageDropzoneWithPreview image={defaultImage} setImage={setImage} defaultValue={defaultImage} />
    )

    let resetImageButton = screen.getByRole('button', { name: 'reset-image' })
    expect(resetImageButton).toBeInTheDocument()
    expect(resetImageButton).toBeDisabled()

    rerender(<LargeImageDropzoneWithPreview image={newImage} setImage={setImage} defaultValue={defaultImage} />)

    resetImageButton = screen.getByRole('button', { name: 'reset-image' })
    expect(resetImageButton).not.toBeDisabled()
    await user.click(resetImageButton)

    expect(setImage).toHaveBeenCalledWith(defaultImage)
  })
})
