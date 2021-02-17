require 'rails_helper'

RSpec.describe Chat, type: :model do
  before do
    @app = App.new(token: 'asddqw823', name: 'myapp')
    @app.save
  end

  it 'chat name cannot be null' do
    chat = Chat.new(app_id: @app.id)
    expect(chat.save).to_not be true
  end

  it 'chat number is unique for apps' do
    chat1 = Chat.create(app_id: @app.id, name: 'mychat', number: 1)
    chat2 = chat1.dup
    expect { chat2.save }.to raise_error(ActiveRecord::RecordNotUnique)
  end
end
