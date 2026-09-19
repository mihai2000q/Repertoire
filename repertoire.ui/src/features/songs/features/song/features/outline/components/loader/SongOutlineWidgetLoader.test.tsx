import { mantineRender } from '../../../../../../../../test-utils.tsx'
import SongOutlineWidgetLoader from './SongOutlineWidgetLoader.tsx'
import { screen } from '@testing-library/react'
import OutlineView from '../../types/enums/OutlineView.ts'

describe('Song Outline Widget Loader', () => {
  it('should render with sections', () => {
    mantineRender(<SongOutlineWidgetLoader outlineView={OutlineView.Sections} />)

    expect(screen.getByLabelText('song-sections-loader')).toBeInTheDocument()
  })

  it('should render with parts', () => {
    mantineRender(<SongOutlineWidgetLoader outlineView={OutlineView.Parts} />)

    expect(screen.getByLabelText('song-parts-loader')).toBeInTheDocument()
  })
})
