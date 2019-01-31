import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class StringsService {

  constructor() { }

  public WebsiteTitle = "Shipwright"

  public Title(word: string): string {
    return word.charAt(0).toUpperCase() + word.slice(1)
  }
}