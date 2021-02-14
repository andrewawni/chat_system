class CreateChats < ActiveRecord::Migration[5.0]
  def change
    create_table :chats do |t|
      t.integer :number
      t.string :name
      t.integer :messages_count, default: 0
      t.belongs_to :app
      t.timestamps
    end
    add_index :chats, :number
  end
end
