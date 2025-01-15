import { reduxRender } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { SongSectionType } from '../../../../types/models/Song.ts'
import { setupServer } from 'msw/node'
import SongSectionTypeSelect from './SongSectionTypeSelect.tsx'

describe('Song Section Type Select', () => {
  const sectionTypes: SongSectionType[] = [
    {
      id: '1',
      name: 'Chorus'
    },
    {
      id: '2',
      name: 'Verse'
    },
    {
      id: '3',
      name: 'Intro'
    }
  ]

  const handlers = [
    http.get('/songs/sections/types', async () => {
      return HttpResponse.json(sectionTypes)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and change types', async () => {
    const user = userEvent.setup()

    const type = sectionTypes[0]
    const newOption = { label: type.name, value: type.id }

    const onChange = vitest.fn()

    const label = 'label'

    const [{ rerender }] = reduxRender(
      <SongSectionTypeSelect label={label} option={null} onChange={onChange} />
    )

    expect(screen.getByRole('textbox', { name: label })).toHaveValue('')

    const select = screen.getByRole('textbox', { name: label })
    await user.click(select)

    for (const st of sectionTypes) {
      expect(await screen.findByRole('option', { name: st.name })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('option', { name: type.name })
    await user.click(selectedOption)

    expect(onChange).toHaveBeenCalledOnce()
    expect(onChange).toHaveBeenCalledWith(newOption)

    rerender(<SongSectionTypeSelect label={label} option={newOption} onChange={onChange} />)

    expect(screen.getByRole('textbox', { name: label })).toHaveValue(type.name)
  })
})
