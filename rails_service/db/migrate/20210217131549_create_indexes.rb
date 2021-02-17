class CreateIndexes < ActiveRecord::Migration[5.0]
  def change
    add_index(:chats, [:app_id, :number], unique: true)
    add_index(:messages, [:chat_id, :number], unique: true)
  end
end
