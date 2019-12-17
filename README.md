# crud
A Golang CRUD application with PostegreSQL - deployed to Heroku. Golang template package for UI.

The UI is simple and straightforward. "HOME" will take you to the table view. "NEW" will take you to the insertion page, where you can add a new record. "TOTALS" will show you the net worth, asset and liability totals. (Net Worth = Assets - Liabilities) 
In the table view, there are "View", "Edit" and "Delete" buttons on the right side - that do the corresponding operations. 

Using PostgreSQL's "bigint" data type and Golang's "int64" type to store dollar values. 
This is because of the lack of a more suitable option (from my limited knowledge of Golang and Postgres) - 
and also the fact that many gateway payment API libraries like Stripe's use int64 and then describe the monetary value in cents 
(eg.: 3500 = $35.00) to keep it simple - which is a fairly common, well understood approach. 

When you enter the "Type" of a new record into the form, please use exactly the words "asset" or "liability" (all lower case) 
because the way the app is set up right now - other values will cause an exception and the insertion will not go through. 

https://morning-headland-60890.herokuapp.com/

