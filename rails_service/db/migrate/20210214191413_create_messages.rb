class CreateMessages < ActiveRecord::Migration[5.0]
  def change
    create_table :messages do |t|
      t.integer :number
      t.text :content
      t.belongs_to :chat
      t.timestamps
    end
    add_index :messages, :number
  end
end
